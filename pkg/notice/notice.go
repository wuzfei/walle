package notice

import "yema.dev/pkg/notice/email"

type Config struct {
	Email email.Config
}

type Msg struct {
}

type Notice interface {
	Send(*Msg) error
}
