package email

import "yema.dev/pkg/notice"

type Config struct {
}

type Email struct {
	*Config
}

func (e *Email) Send(msg *notice.Msg) error {
	return nil
}
