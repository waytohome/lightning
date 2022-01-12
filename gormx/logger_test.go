package gormx

import (
	"errors"
	"testing"
	"time"

	"github.com/waytohome/lightning/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

type GroupUser struct {
	ID        int64     `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"name;size:128;not null;default:'';comment:名字;" json:"name"`
}

func (g GroupUser) TableName() string {
	return "google_group_user"
}

func TestCreateGroupUserTable(t *testing.T) {
	gu := &GroupUser{}
	if DB.Migrator().HasTable(gu) {
		if err := DB.Migrator().DropTable(gu); err != nil {
			t.Fatal(err)
		}
	}
	if err := DB.Migrator().CreateTable(gu); err != nil {
		t.Fatal(err)
	}
}

func TestSelectNotExistRecord(t *testing.T) {
	gu := &GroupUser{}
	err := DB.Model(gu).First(gu, 1).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		t.Fatal(err)
	}
}

func TestInsertRecord(t *testing.T) {
	groupUser := GroupUser{
		Name: "wulalala",
	}
	err := DB.Create(&groupUser).Error
	if err != nil {
		t.Fatal(err)
	}
}

func init() {
	dsn := "root:root@tcp(192.168.31.13:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: WrapLogger(logx.Def, logger.Config{
			SlowThreshold:             10 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	DB = db
}
