package ssh_one

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
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
