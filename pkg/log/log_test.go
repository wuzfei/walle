package log

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewLog(&Config{
		File: "/var/logs//runtime/yema.log",
		//Encoder: "console",
		Level:       "debug",
		Output:      "console",
		Development: true,
	})
	err := errors.New("test err")
	log.Debug("test", zap.Int64("id", 43), zap.Error(err))
}
