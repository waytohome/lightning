package gormx

import (
	"fmt"
	"time"

	"github.com/waytohome/lightning/confx"
	"github.com/waytohome/lightning/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dsnFormatter = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func NewClientWithConfigure(c confx.Configure) (*gorm.DB, error) {
	user, _ := c.GetString("mysql.user", "")
	if user == "" {
		panic("config mysql.user is empty")
	}
	pwd, _ := c.GetString("mysql.password", "")
	host, _ := c.GetString("mysql.host", "")
	if host == "" {
		panic("config mysql.host is empty")
	}
	port, _ := c.GetString("mysql.port", "")
	if port == "" {
		panic("config mysql.port is empty")
	}
	db, _ := c.GetString("mysql.db", "")
	if db == "" {
		panic("config mysql.db is empty")
	}
	dsn := fmt.Sprintf(dsnFormatter, user, pwd, host, port, db)
	return NewClient(dsn)
}

func NewClient(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: WrapLogger(logx.Def, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
	})
}
