package mysql

import "gorm.io/gorm/logger"

type Config struct {
	DSN string

	MaxOpen     int `yaml:"max_open"`
	MaxIdle     int `yaml:"max_idle"`
	MaxLifetime int `yaml:"max_lifetime"`

	LogMode logger.LogLevel `yaml:"log_mode"`
}
