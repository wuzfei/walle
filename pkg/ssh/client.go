package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"sync/atomic"
)

// client ssh方便复用tcp连接管理
type client struct {
	mux          sync.Mutex
	serverConfig *ServerConfig
	client       *ssh.Client
	ref          *int32

	closeFn func()
}

func (s *client) Key() string {
	return s.serverConfig.String()
}

func (s *client) add() {
	atomic.AddInt32(s.ref, 1)
}

func (s *client) done() {
	atomic.AddInt32(s.ref, -1)
}
func (s *client) Ref() int32 {
	return atomic.LoadInt32(s.ref)
}

func (s *client) Done() {
	s.mux.Lock()
	defer s.mux.Unlock()
	n := atomic.AddInt32(s.ref, -1)
	if n == 0 {
		if s.closeFn != nil {
			s.closeFn()
		}
		s.client.Close()
	}
}

func NewClient(conf *ServerConfig, closeFn func()) (_ *client, err error) {
	config := &ssh.ClientConfig{
		User:            conf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(conf.Password), ssh.PublicKeys(conf.IdentitySigner)},
		Timeout:         conf.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	tcpAddress := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	sshClient, err := ssh.Dial("tcp", tcpAddress, config)
	if nil != err {
		return nil, err
	}
	var ref int32
	return &client{
		serverConfig: conf,
		client:       sshClient,
		ref:          &ref,

		closeFn: closeFn,
	}, nil
}

// RunCmd 执行命令
func (s *client) RunCmd(cmd string) (output []byte, err error) {
	s.mux.Lock()
	var session *ssh.Session
	session, err = s.client.NewSession()
	if err != nil {
		s.mux.Unlock()
		return nil, err
	}
	s.add()
	s.mux.Unlock()
	output, err = session.CombinedOutput(cmd)
	_ = session.Close()
	s.Done()
	return
}

func (s *client) NewTerminal(cols, rows int) (term *Terminal, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	session, err := s.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = session.Close()
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", rows, cols, modes)
	if err != nil {
		return
	}
	var reader io.Reader
	var writer io.Writer
	reader, err = session.StdoutPipe()
	if err != nil {
		fmt.Println("session.StdoutPipe error:", err)
		return
	}
	writer, err = session.StdinPipe()
	if err != nil {
		fmt.Println("session.StdinPipe error", err)
		return
	}
	err = session.Shell()
	if err != nil {
		return
	}
	term = &Terminal{
		client:  s,
		session: session,
		reader:  reader,
		writer:  writer,
	}
	s.add()
	return term, nil
}

func (s *client) NewSftp() (*Sftp, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	scp, err := sftp.NewClient(s.client)
	if err != nil {
		return nil, err
	}
	_sftp := &Sftp{
		client:     s,
		sftpClient: scp,
	}
	s.add()
	return _sftp, nil
}

func (s *client) NewRemoteExec(output io.Writer) (*RemoteExec, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_re := &RemoteExec{
		output: output,
		client: s,
	}
	s.add()
	return _re, nil
}
