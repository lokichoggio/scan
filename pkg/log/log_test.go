package log_test

import (
	"testing"

	"scan-eth/pkg/log"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	log.InitLogger(&log.Config{
		File:    "./test.log",
		ErrFile: "./test-err.log",
		Level:   "debug",
		//Encoding: "console",
		Encoding:   "json",
		Stdout:     true,
		MaxSize:    2048,
		MaxBackups: 100,
		MaxAge:     365,
		Compress:   true,
	})

	log.Debugf("test: %s", "debug")
	log.Infof("test: %s", "info")
	log.Warnf("test: %s", "warn")

	// with field
	log.WithFields(zap.String("traceID", "traceID1111")).Error("error")
}
