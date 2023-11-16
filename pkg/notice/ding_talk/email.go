package email

import "yema.dev/pkg/notice"

type Config struct {
}

type DingTalk struct {
	*Config
}

func (e *DingTalk) Send(msg *notice.Msg) error {
	return nil
}
