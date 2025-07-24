package main

import (
	"blog-management/config"
	"blog-management/database"
	"blog-management/internal/routers"
	"blog-management/logger"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	// 1. 读取配置文件
	err := config.InitConfig()
	if err != nil {
		fmt.Printf("config.InitConfig failed. %s \n", err.Error())
		return
	}

	fmt.Println(config.Cfg)

	// 2. 初始化日志文件
	logger.InitLogger(config.Cfg.LogConfig)
	// 3. 初始化数据库
	if !database.InitDb() {
		r := recover()
		if r != nil {
			logger.Logger.Error("failed initdb", zap.String("error", r.(error).Error()))
			return
		}
	}

	// 4. 初始化gin网络框架，初始化路由
	r := routers.InitRouter()
	port := fmt.Sprintf(":%d", config.Cfg.App.Port)
	err = r.Run(port)
	if err != nil {
		logger.Logger.Error("failed start server", zap.String("error", err.Error()))
		return
	}
}
