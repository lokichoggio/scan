package log

import (
	"os"
	"path/filepath"
)

type Config struct {
	File    string `json:"file" yaml:"file"`
	ErrFile string `json:"err_file" yaml:"errFile"`
	// debug info warn error
	Level string `json:"level" yaml:"level"`
	// 日志json格式输出, 如果不输出json, 可以配置成 console
	Encoding string `json:"encoding" yaml:"encoding"`
	// 是否打印到标准输出中
	Stdout bool `json:"stdout" yaml:"stdout"`
	// 日志文件大小限制, 单位: MB
	MaxSize int `json:"max_size" yaml:"maxSize"`
	// 最大保留日志文件数量
	MaxBackups int `json:"max_backups" yaml:"maxBackups"`
	// 日志文件保留天数
	MaxAge int `json:"max_age" yaml:"max_age"`
	// 是否压缩
	Compress bool `json:"compress" yaml:"compress"`
}

func DefaultConfig() *Config {
	path, _ := os.Executable()
	_, execName := filepath.Split(path)

	return &Config{
		File:       execName + ".log",
		ErrFile:    execName + "-err.log",
		Level:      "info",
		Encoding:   "json",
		Stdout:     true,
		MaxSize:    2048,
		MaxBackups: 100,
		MaxAge:     365,
		Compress:   true,
	}
}
