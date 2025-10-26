package main

import (
	"SyncChainEvent/config"
	"SyncChainEvent/logger/xzap"
	service "SyncChainEvent/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("start event sync...")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	onSyncExit := make(chan error, 1)
	go func() {
		defer wg.Done()

		cfg, err := config.UnmarshalConfig("./config.toml") // 读取和解析配置文件
		if err != nil {
			//xzap.WithContext(ctx).Error("Failed to unmarshal config", zap.Error(err))
			fmt.Println("read config err:", err.Error())
			onSyncExit <- err
			return
		}

		_, err = xzap.SetUp(*cfg.Log) // 初始化日志模块
		if err != nil {
			fmt.Println("init log failed, err:", err.Error())
			//xzap.WithContext(ctx).Error("Failed to set up logger", zap.Error(err))
			onSyncExit <- err
			return
		}

		xzap.WithContext(ctx).Info("sync server start", zap.Any("config", cfg))

		s, err := service.New(ctx, cfg) // 初始化服务
		if err != nil {
			xzap.WithContext(ctx).Error("Failed to create sync server", zap.Error(err))
			onSyncExit <- err
			return
		}

		if err := s.Start(); err != nil { // 启动服务
			xzap.WithContext(ctx).Error("Failed to start sync server", zap.Error(err))
			onSyncExit <- err
			return
		}
	}()

	// 信号通知chan
	onSignal := make(chan os.Signal)
	// 优雅退出
	signal.Notify(onSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-onSignal:
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			cancel()
			xzap.WithContext(ctx).Info("Exit by signal", zap.String("signal", sig.String()))
		}
	case err := <-onSyncExit:
		cancel()
		xzap.WithContext(ctx).Error("Exit by error", zap.Error(err))
	}
	wg.Wait()
}
