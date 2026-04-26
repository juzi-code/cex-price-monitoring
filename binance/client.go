package binance

import (
	adshaoBinanceFutures "github.com/adshao/go-binance/v2/futures"
	binanceconnector "github.com/binance/binance-connector-go"
)

var (
	SpotClient    *binanceconnector.Client
	FuturesClient *adshaoBinanceFutures.Client
)

func init() {
	SpotClient = binanceconnector.NewClient("", "")
	FuturesClient = adshaoBinanceFutures.NewClient("", "")
}
