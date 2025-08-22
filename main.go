package main

import (
	binance "cex-price-monitoring/binance"
	binanceWss "cex-price-monitoring/binance/wss"
	"cex-price-monitoring/conf"
	"cex-price-monitoring/constant"
	"cex-price-monitoring/logger"
	"fmt"
	"os"
)

func main() {
	// 初始化日志系统
	config := conf.Cfg()
	if err := logger.Init(config.LogConfig); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("=== CEX价格监控系统启动 ===")

	//订阅现货
	symbols := binance.GetAllSpotUSDTMarket()
	logger.WithField("count", len(symbols)).Info("获取现货市场列表")
	logger.Debugf("现货市场列表: %v", symbols)
	for interval := range constant.PriceChangeThresholdMap {
		var symbolIntervalPair = make(map[string]string)
		for _, symbol := range symbols {
			symbolIntervalPair[symbol] = interval
		}
		go (func() {
			binanceWss.SubKLines(symbolIntervalPair)
		})()
	}

	//订阅期货
	futuresSymbols := binance.GetAllFuturesUSDTMarket()
	logger.WithField("count", len(futuresSymbols)).Info("获取期货市场列表")
	logger.Debugf("期货市场列表: %v", futuresSymbols)
	for interval := range constant.PriceChangeThresholdMap {
		var symbolIntervalPair = make(map[string]string)
		for _, symbol := range futuresSymbols {
			symbolIntervalPair[symbol] = interval
		}
		go (func() {
			binanceWss.SubFuturesKLines(symbolIntervalPair)
		})()
	}
	logger.Info("所有订阅服务已启动，程序进入等待状态...")
	select {}
}
