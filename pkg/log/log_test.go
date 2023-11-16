package log

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

var logConfig = Config{
	File: "/var/logs//runtime/yema.log",
	//Encoder: "console",
	Level:       "debug",
	Output:      "console",
	Development: true,
}

func TestLogger(t *testing.T) {
	log := NewLog(&logConfig)
	err := errors.New("test err")
	log.Debug("test", zap.Int64("id", 43), zap.Error(err))
}
