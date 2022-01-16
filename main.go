package lightning

import (
	"github.com/go-redis/redis/v8"
	"github.com/waytohome/lightning/confx"
	"github.com/waytohome/lightning/ginx"
	"github.com/waytohome/lightning/gormx"
	"github.com/waytohome/lightning/logx"
	"github.com/waytohome/lightning/redix"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB redis.Cmdable
)

func QuickStart() {
	c, err := confx.NewFileConfigure("config.yaml", nil)
	if err != nil {
		panic(err)
	}
	Start(c)
}

func Start(c confx.Configure) {
	level, _ := c.GetString("logger.level", "warn")
	logx.SetLevel(level)
	db, err := gormx.NewClientWithConfigure(c)
	if err != nil {
		panic(err)
	}
	DB = db
	rdb := redix.NewClientWithConfigure(c)
	RDB = rdb
	ginx.InitRoutersWithConfigure(c)
}
