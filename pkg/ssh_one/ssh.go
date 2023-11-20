package ssh_one

import (
	"context"
	"fmt"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"time"
)

var (
	ErrSSH = errs.Class("ssh")
)

// Config ssh配置
type Config struct {
	IdentityFile     string        `help:"免密登陆密钥地址" default:"$HOME/.ssh/id_rsa"`
	IdentityPassword string        `help:"免密登陆密钥密码" default:""`
	Timeout          time.Duration `help:"连接超时" default:"30s"`
}

type Server struct {
	Host           string        `json:"host"`
	User           string        `json:"user"`
	Password       string        `json:"password"` //如果密码为空，则认为是免密登陆
	Port           int           `json:"port"`
	Timeout        time.Duration `json:"timeout"`
	IdentitySigner ssh.Signer    `json:"-"`
}

func (conf *Config) IdentitySigner() (signer ssh.Signer, err error) {
	if conf == nil {
		return nil, ErrSSH.New("config can't nil ")
	}
	_, err = os.Stat(conf.IdentityFile)
	if err != nil {
		return nil, ErrSSH.New("ssh config IdentityFile: %s not exists", conf.IdentityFile)
	}
	bytes, err := os.ReadFile(conf.IdentityFile)
	if err != nil {
		return nil, ErrSSH.Wrap(err)
	}
	signer, err = ssh.ParsePrivateKey(bytes)
	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(bytes, []byte(conf.IdentityPassword))
	}
	if err != nil {
		err = ErrSSH.Wrap(err)
	}
	return
}

func (s *Server) ID() (key string) {
	return fmt.Sprintf("%s:%s@%s:%d", s.User, s.Password, s.Host, s.Port)
}

type item struct {
	client   *client
	lastTime time.Time
}

type Ssh struct {
	mux   *sync.RWMutex
	items map[string]*item
	//pool *pool
	identitySigner ssh.Signer
	timeout        time.Duration

	log *zap.Logger
}

func NewSsh(ctx context.Context, conf *Config, log *zap.Logger) (*Ssh, error) {
	iSigner, err := conf.IdentitySigner()
	if err != nil {
		return nil, err
	}
	sh := &Ssh{
		mux:   &sync.RWMutex{},
		items: make(map[string]*item),

		identitySigner: iSigner,
		timeout:        conf.Timeout,

		log: log.Named("ssh"),
	}
	go sh.start(ctx)
	return sh, nil
}

func (s *Ssh) start(ctx context.Context) {
	tk := time.NewTicker(time.Second * 5)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			s.mux.Lock()
			if len(s.items) == 0 {
				s.log.Debug("当前clients:[0]")
			} else {
				t := time.Now()
				for k, v := range s.items {
					flag := false
					s.log.Debug(fmt.Sprintf("当前clients:[%s], sessions:[%d], lastTime:[%s]", k, v.client.Ref(), v.lastTime))
					if v.client.Ref() == 0 {
						if !v.lastTime.Add(time.Second * 10).After(t) {
							flag = true
						}
					} else {
						if !v.lastTime.Add(time.Second * 30).After(t) {
							flag = true
						}
					}
					if flag {
						v.client.Close()
						delete(s.items, k)
					}
				}
			}
			s.mux.Unlock()
		}
	}
}

func (s *Ssh) getClient(server *Server) (*client, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if item, ok := s.items[server.ID()]; ok {
		s.items[server.ID()].lastTime = time.Now()
		return item.client, nil
	}
	c, err := newClient(server)
	if err != nil {
		return nil, err
	}
	s.items[server.ID()] = &item{
		client:   c,
		lastTime: time.Now(),
	}
	return c, err
}

// NewTerminal 获取会话终端
func (s *Ssh) NewTerminal(server *Server, cols, rows int) (sess *Terminal, err error) {
	sshClient, err := s.getClient(server)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return
	}
	return sshClient.NewTerminal(cols, rows)
}

func (s *Ssh) NewSftp(server *Server) (*Sftp, error) {
	sshClient, err := s.getClient(server)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return nil, err
	}
	return sshClient.NewSftp()
}
