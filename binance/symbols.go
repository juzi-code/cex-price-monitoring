package binance

import (
	"context"
	"strings"

	adshaoBinanceFutures "github.com/adshao/go-binance/v2/futures"
	binanceconnector "github.com/binance/binance-connector-go"

	"cex-price-monitoring/logger"
)

// Symbol holds the symbol name and its base asset.
type Symbol struct {
	Name      string
	BaseAsset string
}

func GetSpotUSDTSymbols() []Symbol {
	res, err := SpotClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		logger.WithField("error", err).Error("获取现货交易对信息失败")
		return nil
	}
	return filterSpotSymbols(res.Symbols)
}

func GetFuturesUSDTSymbols() []Symbol {
	res, err := FuturesClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		logger.WithField("error", err).Error("获取期货交易对信息失败")
		return nil
	}
	return filterFuturesSymbols(res.Symbols)
}

func GetSpotTicker24h(symbol string) *binanceconnector.Ticker24hrResponse {
	arr, err := SpotClient.NewTicker24hrService().Symbol(symbol).Do(context.Background())
	if err != nil {
		logger.WithFields(logger.Fields{"symbol": symbol, "error": err}).Error("获取现货24h行情失败")
		return nil
	}
	if len(arr) == 0 {
		return nil
	}
	return arr[0]
}

func GetFuturesTicker24h(symbol string) *adshaoBinanceFutures.PriceChangeStats {
	arr, err := FuturesClient.NewListPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		logger.WithFields(logger.Fields{"symbol": symbol, "error": err}).Error("获取期货24h行情失败")
		return nil
	}
	if len(arr) == 0 {
		return nil
	}
	return arr[0]
}

func filterSpotSymbols(raw []*binanceconnector.SymbolInfo) []Symbol {
	out := make([]Symbol, 0, len(raw))
	for _, s := range raw {
		if strings.HasSuffix(s.Symbol, "USDT") {
			out = append(out, Symbol{Name: s.Symbol, BaseAsset: s.BaseAsset})
		}
	}
	return out
}

func filterFuturesSymbols(raw []adshaoBinanceFutures.Symbol) []Symbol {
	out := make([]Symbol, 0, len(raw))
	for _, s := range raw {
		if strings.HasSuffix(s.Symbol, "USDT") {
			out = append(out, Symbol{Name: s.Symbol, BaseAsset: s.BaseAsset})
		}
	}
	return out
}
