package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func InitConn(c *Config) (*DB, error) {
	d := mysql.New(mysql.Config{
		DSN: c.DSN,
	})

	db, err := gorm.Open(d, &gorm.Config{
		Logger: logger.Default.LogMode(c.LogMode),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(c.MaxOpen)
	sqlDB.SetMaxIdleConns(c.MaxIdle)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxLifetime) * time.Second)

	return &DB{db}, nil
}

func (d *DB) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
