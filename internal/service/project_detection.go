package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"yema.dev/api"
	"yema.dev/internal/model"
	"yema.dev/pkg/repo"
	"yema.dev/pkg/ssh"
)

const (
	msgTypeMsg     = 0
	msgTypeSuccess = 1
	msgTypeError   = 2
)

type write struct {
	mux sync.Mutex
	buf bytes.Buffer
}

func (w *write) Write(b []byte) (n int, err error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.buf.Write(b)
}

func (w *write) Reset() {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.buf.Reset()
}
func (w *write) String() string {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.buf.String()
}

type DetectionMsg struct {
	ServerId int64  `json:"server_id"`
	Title    string `json:"title"`
	Error    string `json:"error"`
	Todo     string `json:"todo"`
}

func sendDetectionMsgFn(ch chan<- *DetectionMsg) func(title, todo, err string, serverId int64) {
	return func(title, todo, err string, serverId int64) {
		ch <- &DetectionMsg{
			ServerId: serverId,
			Title:    title,
			Error:    err,
			Todo:     todo,
		}
	}
}

// DetectionWs 项目检测
func (s *projectService) DetectionWs(ctx context.Context, wsConn *websocket.Conn, spaceWithId *api.SpaceWithId) (err error) {
	if s.detectionTimeout == 0 {
		s.detectionTimeout = time.Second * 30
	}
	ctx, cancel := context.WithTimeout(ctx, s.detectionTimeout)
	dMsgChan := make(chan *DetectionMsg)
	defer func() {
		close(dMsgChan)
	}()
	sendMsg := sendDetectionMsgFn(dMsgChan)

	//检测逻辑
	go func() {
		var err error
		defer func() {
			if _err := recover(); _err != nil {
				s.log.Error("1.DetectionWs 已关闭写入渠道", zap.Any("panic", _err))
			}
			cancel()
		}()

		s.log.Info("获取数据库项目信息")
		project, err := s.Detail(ctx, spaceWithId)
		if err != nil {
			sendMsg("检测项目不存在", "请检查项目是否存在，或者刷新页面再尝试", "", 0)
			return
		}

		s.log.Info("clone仓库代码")
		_, err = s.repo.New(repo.TypeRepo(project.RepoType), project.RepoUrl, strconv.Itoa(int(project.ID)))
		if err != nil {
			sendMsg("代码clone失败", "1、请检查仓库地址："+project.RepoUrl+"是否正确；\n 2、请检查"+project.RepoType+"相关配置是否正确", err.Error(), 0)
			return
		}
		if len(project.Servers) == 0 {
			sendMsg("项目未绑定发布服务器", "请添加发布服务器后，在修改项目重新选择绑定", "", 0)
			err = errors.New("项目未绑定发布服务器")
			return
		}

		s.log.Info("检查服务器配置信息")
		g := sync.WaitGroup{}
		//检查服务器
		for _, server := range project.Servers {
			g.Add(1)
			go func(server model.Server) {
				var err error
				defer func() {
					if _err := recover(); _err != nil {
						s.log.Error("2.DetectionWs 已关闭写入渠道", zap.Any("panic", _err))
					}
					if err != nil {
						s.log.Error("项目检测服务器失败",
							zap.Int64("ProjectId", project.ID),
							zap.String("server", server.Hostname()),
							zap.Error(err))
					}
					g.Done()
				}()

				s.log.Info("连接服务器", zap.String("server", server.Hostname()))
				buf := write{}
				var re *ssh.RemoteExec
				re, err = s.ssh.NewRemoteExec(ssh.ServerConfig{User: server.User, Port: server.Port, Host: server.Host}, &buf)
				if err != nil {
					sendMsg("远程目标机器免密码登录连接失败",
						fmt.Sprintf("在宿主机中配置免密码登录，把宿主机用户[%s]的~/.ssh/id_rsa.pub添加到远程目标机器用户[%s]的~/.ssh/authorized_keys", server.User, server.User),
						err.Error(),
						server.ID)
					return
				}
				defer func() { _ = re.Close() }()

				s.log.Info("服务器创建发布目录", zap.String("server", server.Hostname()))
				webroot := filepath.Dir(project.TargetRoot)
				cmd := fmt.Sprintf("[ -d %s ] || mkdir -p %s", webroot, webroot)
				err = re.Run(cmd)
				if err != nil {
					sendMsg("远程目标机器创建目录失败",
						fmt.Sprintf("请检查远程目标服务器用户[%s]的权限，失败执行命令：%s", server.User, cmd),
						err.Error(),
						server.ID)
					return
				}

				s.log.Info("服务器检查软链接", zap.String("server", server.Hostname()))
				cmd = fmt.Sprintf("[ -L \"%s\" ] && echo \"true\" || echo \"false\"", project.TargetRoot)
				buf.Reset()
				err = re.Run(cmd)
				if err != nil {
					sendMsg("目标机器执行命令失败",
						fmt.Sprintf("请检查远程目标服务器用户[%s]的权限，失败执行命令：%s", server.User, cmd),
						err.Error(),
						server.ID)
					return
				}
				if buf.String() == "false" {
					err = fmt.Errorf("远程目标机器%s webroot不能是已存在的目录，必须为软链接，你不必新建，walle会自行创建。", server.Host)
					sendMsg("远程目标机器webroot不能是已建好的目录",
						"手工删除远程目标机器："+server.Host+" webroot目录："+project.TargetRoot,
						err.Error(),
						server.ID)
					return
				}
			}(server)

		}
		g.Wait()
	}()

	//客户端发送消息
	var res []byte
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-dMsgChan:
			res, err = json.Marshal(msg)
			s.log.Debug("项目检测回送客户端消息", zap.ByteString("msg", res), zap.Error(err))
			if err != nil {
				s.log.Error("DetectionWs json.Marshal error", zap.Error(err))
				continue
			}
			err = wsConn.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				cancel()
				return fmt.Errorf("DetectionWs wsConn.WriteMessage error:%s", err)
			}
		default:
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				s.log.Error("ws.ReadMessage err:", zap.ByteString("msg", msg), zap.Error(err))
				cancel()
				return nil
			}
		}
	}
}
