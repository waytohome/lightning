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
	maxConnect, _ := c.GetInt("mysql.connect.max", 100)
	idleConnect, _ := c.GetInt("mysql.connect.idle", 10)
	maxConnectLife, _ := c.GetInt("mysql.connect.max-life", 60)
	idleConnectLife, _ := c.GetInt("mysql.connect-idle-life", 60)
	return NewClient(Config{
		Dsn:             dsn,
		MacConnect:      maxConnect,
		IdleConnect:     idleConnect,
		ConnectMaxLife:  time.Duration(maxConnectLife) * time.Minute,
		ConnectIdleLife: time.Duration(idleConnectLife) * time.Minute,
	})
}

type Config struct {
	Dsn             string
	MacConnect      int
	IdleConnect     int
	ConnectMaxLife  time.Duration
	ConnectIdleLife time.Duration
}

func NewClient(conf Config) (*gorm.DB, error) {
	logLv := logger.Silent
	lv := logx.GetLevel()
	for k, v := range levelMap {
		if v == lv {
			logLv = k
		}
	}
	config := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  logLv,
	}
	db, err := gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{
		Logger: WrapLogger(logx.Def, config),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(conf.MacConnect)
	sqlDB.SetMaxIdleConns(conf.IdleConnect)
	sqlDB.SetConnMaxLifetime(conf.ConnectMaxLife)
	sqlDB.SetConnMaxIdleTime(conf.ConnectIdleLife)
	return db, err
}
