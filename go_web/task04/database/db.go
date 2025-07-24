package database

import (
	"blog-management/database/model"
	"blog-management/logger"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func InitDb() bool {
	var err error
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "file:blog.db?cache=shared",
	}, &gorm.Config{Logger: NewZapGormLogger(logger.Logger)})
	if err != nil {
		fmt.Printf("gorm.Open failed. %s \n", err.Error())
		panic("failed to connect database")
		return false
	}

	// 自动迁移模型
	_ = DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	return true
}

func NewZapGormLogger(log *zap.Logger) *ZapGormLogger {
	return &ZapGormLogger{log: log}
}
