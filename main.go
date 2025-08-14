package main

import (
	binance "cex-price-monitoring/binance"
	binanceWss "cex-price-monitoring/binance/wss"
	"cex-price-monitoring/constant"
	"fmt"
)

func main() {
	//订阅现货
	symbols := binance.GetAllSpotUSDTMarket()
	fmt.Println("现货-所有市场:", symbols)
	fmt.Println("现货-市场数量:", len(symbols))
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
	fmt.Println("期货-所有市场:", futuresSymbols)
	fmt.Println("期货-市场数量:", len(futuresSymbols))
	for interval := range constant.PriceChangeThresholdMap {
		var symbolIntervalPair = make(map[string]string)
		for _, symbol := range symbols {
			symbolIntervalPair[symbol] = interval
		}
		go (func() {
			binanceWss.SubFuturesKLines(symbolIntervalPair)
		})()
	}
	select {}
}
