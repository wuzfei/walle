package service

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"time"
	"unicode/utf8"
	"yema.dev/api"
	"yema.dev/api/server"
	"yema.dev/internal/errcode"
	"yema.dev/internal/model"
	"yema.dev/internal/model/field"
	"yema.dev/internal/repository"
	"yema.dev/pkg/ssh"
)

type ServerService interface {
	List(ctx context.Context, req *server.ListReq) (total int64, res []*model.Server, err error)
	Create(ctx context.Context, req *server.CreateReq) error
	Update(ctx context.Context, req *server.UpdateReq) error
	Delete(ctx context.Context, req *api.SpaceWithId) error

	Check(ctx context.Context, req *api.SpaceWithId) error
	SetAuthorized(ctx context.Context, req *server.SetAuthorizedReq) error
	Terminal(ctx context.Context, wsConn *websocket.Conn, spaceWithId *api.SpaceWithId, username string) error
}

func NewServerService(service *Service, serverRepo repository.ServerRepository, ssh *ssh.Ssh) ServerService {
	return &serverService{
		serverRepo: serverRepo,
		Service:    service,
		ssh:        ssh,
	}
}

type serverService struct {
	serverRepo repository.ServerRepository
	ssh        *ssh.Ssh
	*Service
}

func (s *serverService) List(ctx context.Context, req *server.ListReq) (total int64, list []*model.Server, err error) {
	return s.serverRepo.List(ctx, req)
}

func (s *serverService) Create(ctx context.Context, req *server.CreateReq) error {
	m := &model.Server{
		SpaceId:     req.SpaceId,
		Name:        req.Name,
		Host:        req.Host,
		Port:        req.Port,
		User:        req.User,
		Status:      field.StatusDisable,
		Description: req.Description,
	}
	_m, err := s.serverRepo.FindByHostIp(ctx, m.SpaceId, m.User, m.Host, m.Port)
	if err != nil {
		return err
	}
	if _m.ID != 0 {
		return errors.New(fmt.Sprintf("已存在该主机：[%s@%s:%d]", m.User, m.Host, m.Port))
	}
	return s.serverRepo.Create(ctx, m)
}

func (s *serverService) Update(ctx context.Context, req *server.UpdateReq) error {
	m, err := s.serverRepo.FindByHostIp(ctx, req.SpaceId, req.User, req.Host, req.Port)
	if err != nil {
		return err
	}
	if m.ID == 0 || m.ID != req.ID {
		return errors.New("更新错误")
	}
	m.Name = req.Name
	m.User = req.User
	m.Host = req.Host
	m.Port = req.Port
	m.Description = req.Description
	return s.serverRepo.UpdateFields(ctx, &m, req.Fields()...)
}

// Delete 删除
func (s *serverService) Delete(ctx context.Context, req *api.SpaceWithId) error {
	return s.tm.Transaction(ctx, func(ctx context.Context) error {
		m, err := s.serverRepo.GetByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if m.SpaceId != req.SpaceId {
			return errcode.ErrBadRequest
		}
		if err := s.serverRepo.ClearProjects(ctx, req.ID); err != nil {
			return err
		}
		return s.serverRepo.DeleteByID(ctx, m.ID)
	})
}

const (
	waringMsg = iota
	errorMsg
	successMsg
	defaultMsg
)

const (
	connectTimeout     = time.Minute * 10 //保持连接最长时间
	buffTime           = time.Microsecond * 500
	wsMsgTypeResize    = "resize"
	wsMsgTypeCmd       = "cmd"
	wsMsgTypeHeartbeat = "ping"
)

type TerminalWsMsg struct {
	Typ string `json:"typ"`
	Cmd string `json:"cmd,omitempty"`
	Col int    `json:"col,omitempty"`
	Row int    `json:"row,omitempty"`
}

func (s *serverService) Check(ctx context.Context, req *api.SpaceWithId) error {
	m, err := s.serverRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}

	output, err := s.ssh.RunCmd(ssh.ServerConfig{
		User: m.User,
		Host: m.Host,
		Port: m.Port,
	}, "pwd")
	s.log.Debug("CheckConnect", zap.String("cmd", "pwd"), zap.ByteString("output", output), zap.Error(err))
	if err != nil && m.Status.IsEnable() {
		m.Status = field.StatusDisable
		return s.serverRepo.UpdateFields(ctx, m, "status")
	}
	if err == nil && m.Status.IsDisable() {
		m.Status = field.StatusEnable
		return s.serverRepo.UpdateFields(ctx, m, "status")
	}
	return err
}

func (s *serverService) SetAuthorized(ctx context.Context, req *server.SetAuthorizedReq) error {
	m, err := s.serverRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if m.SpaceId != req.SpaceId {
		return errcode.ErrBadRequest
	}

	if m.Status.IsEnable() {
		return errors.New("该服务器能正常连接，无需设置")
	}
	signer := s.ssh.GetIdentitySigner()
	hostname, _ := os.Hostname()
	publicKeyStr := fmt.Sprintf("%s %s %s", signer.PublicKey().Type(), base64.StdEncoding.EncodeToString(signer.PublicKey().Marshal()), hostname)
	runCmd := fmt.Sprintf("mkdir -p $HOME/.ssh && echo '%s' >> $HOME/.ssh/authorized_keys && chmod 600 $HOME/.ssh/authorized_keys", publicKeyStr)
	output, err := s.ssh.RunCmd(ssh.ServerConfig{
		User:     m.User,
		Host:     m.Host,
		Password: req.Password,
		Port:     m.Port,
	}, runCmd)
	s.log.Debug("Setting", zap.String("cmd", runCmd), zap.ByteString("output", output), zap.Error(err))
	if err == nil {
		m.Status = field.StatusEnable
		if _err := s.serverRepo.UpdateFields(ctx, m, "status"); _err != nil {
			s.log.Error("更新数据库失败", zap.Int64("server_id", m.ID), zap.Int("status", field.StatusEnable))
		}
	}
	return err
}

func (s *serverService) Terminal(ctx context.Context, wsConn *websocket.Conn, spaceWithId *api.SpaceWithId, username string) error {
	wsSendMsg := func(msg string, msgType int) error {
		_err := wsConn.WriteMessage(websocket.TextMessage, []byte(terminalMsg(msg, msgType)))
		if _err != nil {
			s.log.Debug("发送ws数据失败", zap.Error(_err))
		}
		return _err
	}

	m, err := s.serverRepo.GetByID(ctx, spaceWithId.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = wsSendMsg("该服务器不存在！", errorMsg)
		} else {
			_ = wsSendMsg(err.Error(), errorMsg)
		}
		return err
	}
	if m.SpaceId != spaceWithId.SpaceId {
		return errcode.ErrBadRequest
	}

	if err = wsSendMsg("正在连接服务器...", successMsg); err != nil {
		return err
	}
	sshTerminal, err := s.ssh.NewTerminal(ssh.ServerConfig{
		User:     m.User,
		Host:     m.Host,
		Password: "",
		Port:     m.Port,
	}, 200, 40)
	if err != nil {
		_ = wsSendMsg(err.Error(), errorMsg)
		return err
	}
	defer func() {
		_ = sshTerminal.Close()
	}()
	if err = wsSendMsg("连接服务器成功！", successMsg); err != nil {
		return err
	}
	if err = wsSendMsg("Hello "+username+"，您所操作的所有命令都将会被记录，请谨慎操作！！！", waringMsg); err != nil {
		return err
	}
	s.dealMsg(ctx, wsConn, sshTerminal)
	return nil
}

// dealMsg 终端数据交互
func (s *serverService) dealMsg(ctx context.Context, wsConn *websocket.Conn, sshTerminal *ssh.Terminal) {
	connectTimeoutT := time.NewTimer(connectTimeout)
	bufTimeT := time.NewTimer(buffTime)
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		connectTimeoutT.Stop()
		bufTimeT.Stop()
	}()

	go func() {
		var err error
		defer func() {
			if err != nil {
				s.log.Error("处理ws消息出错", zap.Error(err))
			}
		}()
		var msg []byte
		for {
			select {
			//监听上下文退出
			case <-ctx.Done():
				return
			default:
				_, msg, err = wsConn.ReadMessage()
				if err != nil {
					cancel()
					s.log.Error("ws.ReadMessage err:", zap.Error(err))
					return
				}
				wsMsg := new(TerminalWsMsg)
				if err = json.Unmarshal(msg, wsMsg); err != nil {
					continue
				}
				switch wsMsg.Typ {
				case wsMsgTypeResize:
					err = sshTerminal.WindowChange(wsMsg.Row, wsMsg.Col)
				case wsMsgTypeHeartbeat:
					wsConn.WriteMessage(websocket.TextMessage, []byte("pong"))
				default:
					_, err = sshTerminal.Write([]byte(wsMsg.Cmd))
				}
				if err != nil {
					cancel()
					s.log.Error("sshTerminal command err:", zap.Error(err))
					return
				}
			}
		}
	}()

	r := make(chan rune)
	//读取buf
	br := bufio.NewReader(sshTerminal)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				x, size, err := br.ReadRune()
				if err != nil {
					cancel()
					s.log.Error("读取终端消息出错", zap.Error(err))
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}
	}()

	buf := make([]byte, 0)
	// 主循环
	for {
		select {
		case <-connectTimeoutT.C:
			cancel()
			return
		case <-ctx.Done():
			return
		case <-bufTimeT.C:
			if len(buf) != 0 {
				err := wsConn.WriteMessage(websocket.TextMessage, buf)
				buf = []byte{}
				if err != nil {
					cancel()
					s.log.Error("ws.WriteMessage err:", zap.Error(err))
					return
				}
				connectTimeoutT.Reset(connectTimeout)
			}
			bufTimeT.Reset(buffTime)
		case d := <-r:
			if d != utf8.RuneError {
				p := make([]byte, utf8.RuneLen(d))
				utf8.EncodeRune(p, d)
				buf = append(buf, p...)
			} else {
				buf = append(buf, []byte("@")...)
			}
			connectTimeoutT.Reset(connectTimeout)
		}
	}
}

func terminalMsg(msg string, typ int) string {
	switch typ {
	case waringMsg:
		return fmt.Sprintf("\x1b[33m%s\x1b[m\r\n", msg)
	case errorMsg:
		return fmt.Sprintf("\x1b[31m%s\x1b[m\r\n", msg)
	case successMsg:
		return fmt.Sprintf("\x1b[32m%s\x1b[m\r\n", msg)
	default:
		return fmt.Sprintf("%s\r\n", msg)
	}
}
