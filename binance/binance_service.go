package binance

import (
	"context"
	"fmt"
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		fmt.Println(err)
		return nil
	}
	if len(ticker24hrResponseArr) == 0 {
		fmt.Println("getTicker24h result null")
		return nil
	}
	return ticker24hrResponseArr[0]
}
func GetTicker24hFutures(symbol string) *adshaoBinance.PriceChangeStats {
	ticker24hrResponseArr, err := clientFutures.NewListPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(ticker24hrResponseArr) == 0 {
		fmt.Println("etTicker24hFutures result null")
		return nil
	}
	return ticker24hrResponseArr[0]
}
