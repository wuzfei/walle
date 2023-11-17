package ssh_one

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
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
}

func NewSsh(conf *Config) *Ssh {
	return &Ssh{
		//pool: &pool{
		//	items: make(map[string]*Client),
		//	newFn: newClient(conf),
		//},
	}
}

func (s *Ssh) getClient(server *Server) (*client, error) {
	s.mux.Lock()
	defer s.mux.Lock()
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

func (s *Ssh) NewSftp(server *Server) (*Sftp, error) {
	s.mux.RLock()
	if item, ok := s.items[server.ID()]; ok {

	}
}
