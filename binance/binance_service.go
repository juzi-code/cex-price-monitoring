package binance

import (
	"cex-price-monitoring/logger"
	"context"
	adshaoBinance "github.com/adshao/go-binance/v2/futures"
	binanceconnector "github.com/binance/binance-connector-go"
	"strings"
)

var client *binanceconnector.Client
var clientFutures *adshaoBinance.Client

func init() {
	client = binanceconnector.NewClient("", "")
	clientFutures = adshaoBinance.NewClient("", "")
}

// GetAllSpotUSDTMarket 获取所有USDT现货交易对
func GetAllSpotUSDTMarket() []string {
	allExchange := client.NewExchangeInfoService()
	res, err := allExchange.Do(context.Background())
	if err != nil {
		logger.WithField("error", err).Error("获取现货交易对信息失败")
		return nil
	}
	//打印res.Symbols各个元素的值
	//初始化一个string数组
	var symbols []string
	//res.Symbols装进symbols
	for _, v := range res.Symbols {
		if strings.HasSuffix(v.Symbol, "USDT") {
			//数组追加元素
			symbols = append(symbols, v.Symbol)
		}
	}
	return symbols
}

// GetAllSpotUSDTMarket 获取所有USDT期货交易对
func GetAllFuturesUSDTMarket() []string {
	allExchange := clientFutures.NewExchangeInfoService()
	res, err := allExchange.Do(context.Background())
	if err != nil {
		logger.WithField("error", err).Error("获取期货交易对信息失败")
		return nil
	}
	//打印res.Symbols各个元素的值
	//初始化一个string数组
	var symbols []string
	//res.Symbols装进symbols
	for _, v := range res.Symbols {
		if strings.HasSuffix(v.Symbol, "USDT") {
			//数组追加元素
			symbols = append(symbols, v.Symbol)
		}
	}
	return symbols
}

func GetTicker24h(symbol string) *binanceconnector.Ticker24hrResponse {
	ticker24hrResponseArr, err := client.NewTicker24hrService().Symbol(symbol).Do(context.Background())
	if err != nil {
		logger.WithFields(logger.Fields{"symbol": symbol, "error": err}).Error("获取现货24h行情失败")
		return nil
	}
	if len(ticker24hrResponseArr) == 0 {
		logger.WithField("symbol", symbol).Warn("现货24h行情数据为空")
		return nil
	}
	return ticker24hrResponseArr[0]
}
func GetTicker24hFutures(symbol string) *adshaoBinance.PriceChangeStats {
	ticker24hrResponseArr, err := clientFutures.NewListPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		logger.WithFields(logger.Fields{"symbol": symbol, "error": err}).Error("获取期货24h行情失败")
		return nil
	}
	if len(ticker24hrResponseArr) == 0 {
		logger.WithField("symbol", symbol).Warn("期货24h行情数据为空")
		return nil
	}
	return ticker24hrResponseArr[0]
}
