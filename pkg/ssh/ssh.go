package ssh

import (
	"fmt"
	"github.com/zeebo/errs"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"time"
	"yema.dev/app/utils"
)

var (
	ErrSSH = errs.Class("ssh")
)

type Ssh struct {
	mux     *sync.RWMutex
	clients map[string]*client

	identitySigner ssh.Signer
	timeout        time.Duration
}

func NewSSH(conf *Config) (*Ssh, error) {
	iSigner, err := conf.IdentitySigner()
	if err != nil {
		return nil, err
	}
	sh := &Ssh{
		mux:     &sync.RWMutex{},
		clients: make(map[string]*client),

		identitySigner: iSigner,
		timeout:        conf.Timeout,
	}

	if utils.IsDev() {
		go func() {
			tk := time.NewTicker(time.Second * 5)
			defer tk.Stop()
			for {
				select {
				case <-tk.C:
					sh.mux.RLock()
					if len(sh.clients) == 0 {
						fmt.Printf("当前clients:[0] \r\n")
					} else {
						for _, v := range sh.clients {
							fmt.Printf("当前clients:[%s], sessions:[%d]\r\n", v.serverConfig.String(), v.Ref())
						}
					}
					sh.mux.RUnlock()
				}
			}
		}()
	}
	return sh, nil
}

func (s *Ssh) newClient(conf ServerConfig) (sc *client, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := conf.String()
	if v, ok := s.clients[key]; ok {
		return v, nil
	}
	if conf.Timeout == 0 {
		conf.Timeout = s.timeout
	}
	conf.IdentitySigner = s.identitySigner
	sc, err = NewClient(&conf, func() {
		s.mux.Lock()
		defer s.mux.Unlock()
		key := conf.String()
		if _, ok := s.clients[key]; ok {
			delete(s.clients, key)
		}
	})
	if err != nil {
		return
	}
	s.clients[key] = sc
	return
}

func (s *Ssh) removeClient(conf *ServerConfig) {
	s.mux.Lock()
	defer s.mux.Unlock()
	key := conf.String()
	if _, ok := s.clients[key]; ok {
		delete(s.clients, key)
	}
}

// NewTerminal 获取会话终端
func (s *Ssh) NewTerminal(conf ServerConfig, cols, rows int) (sess *Terminal, err error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return
	}
	return sshClient.NewTerminal(cols, rows)
}

// RunCmd 直接连接执行命令
func (s *Ssh) RunCmd(conf ServerConfig, cmd string) (output []byte, err error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return
	}
	return sshClient.RunCmd(cmd)
}

func (s *Ssh) NewSftp(conf ServerConfig) (*Sftp, error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return nil, err
	}
	return sshClient.NewSftp()
}

func (s *Ssh) NewRemoteExec(conf ServerConfig, output io.Writer) (*RemoteExec, error) {
	sshClient, err := s.newClient(conf)
	if err != nil {
		err = ErrSSH.Wrap(err)
		return nil, err
	}
	return sshClient.NewRemoteExec(output)
}

func (s *Ssh) GetIdentitySigner() ssh.Signer {
	return s.identitySigner
}
