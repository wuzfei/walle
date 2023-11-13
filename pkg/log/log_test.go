package log

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewLog(&Config{
		File: "/var/logs//runtime/yema.log",
		//Encoder: "console",
		Level:  "debug",
		Output: "all",
	})
	log.Debug("test", zap.Int64("id", 43))
}
