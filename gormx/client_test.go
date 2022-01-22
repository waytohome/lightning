package gormx

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/waytohome/lightning/logx"
	"gorm.io/gorm"
)

type Demo struct {
	gorm.Model
}

func TestConnectPool(t *testing.T) {
	db, err := NewClient(Config{
		Dsn:             "root:root@tcp(192.168.31.13:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		MacConnect:      100,
		IdleConnect:     5,
		ConnectMaxLife:  time.Hour,
		ConnectIdleLife: time.Minute * 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	demo := Demo{}
	if !db.Migrator().HasTable(&demo) {
		if err := db.Migrator().CreateTable(&demo); err != nil {
			t.Fatal(err)
		}
	}
	pool, _ := ants.NewPool(100)
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		pool.Release()
	}()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		_ = pool.Submit(func() {
			defer wg.Done()
			demo := Demo{}
			err := db.First(&demo, 1).Error
			if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
				logx.Error("select err", logx.String("err", err.Error()))
			}
		})
	}
}

func init() {
	logx.SetLevel("warn")
}
