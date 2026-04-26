package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cex-price-monitoring/conf"
	"cex-price-monitoring/logger"
	"cex-price-monitoring/tgbot"
	"cex-price-monitoring/binance/wss"
)

func main() {
	cfg := conf.Cfg()

	if err := logger.Init(cfg.LogConfig); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}
	logger.WithField("env", env).Info("=== CEX价格监控系统启动 ===")

	if err := tgbot.Init(cfg.TelegramData.BotToken); err != nil {
		logger.Fatalf("初始化Telegram失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spotManager := wss.NewManager(wss.Spot)
	futuresManager := wss.NewManager(wss.Futures)

	go spotManager.Run(ctx)
	go futuresManager.Run(ctx)

	logger.Info("所有订阅服务已启动，等待信号退出...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("收到退出信号，正在关闭...")
	cancel()
}
