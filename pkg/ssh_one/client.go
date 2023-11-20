package ssh_one

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
	"sync/atomic"
)

type client struct {
	mux    sync.Mutex
	server *Server
	client *ssh.Client
	ref    *int32
}

func newClient(server *Server) (_ *client, err error) {
	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            []ssh.AuthMethod{ssh.Password(server.Password), ssh.PublicKeys(server.IdentitySigner)},
		Timeout:         server.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	tcpAddress := fmt.Sprintf("%s:%d", server.Host, server.Port)
	sshClient, err := ssh.Dial("tcp", tcpAddress, config)
	if nil != err {
		return nil, err
	}
	var ref int32
	return &client{
		server: server,
		client: sshClient,
		ref:    &ref,
	}, nil
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
	atomic.AddInt32(s.ref, -1)
}

func (s *client) Close() error {
	if atomic.LoadInt32(s.ref) == 0 {

	} else {

	}
	return nil
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
