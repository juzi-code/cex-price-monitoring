package data

import "time"

type PriceChangeSignal struct {
	Type               string
	Symbol             string
	CexName            string
	Interval           string
	OpenPrice          float64
	ClosePrice         float64
	HighPrice          float64
	LowPrice           float64
	QuoteVolume        float64
	PriceChangePercent float64
	TradeNum           int64
	Time               time.Time
	CoinTicker         *CoinTicker
}

type CoinTicker struct {
	Symbol             string
	OpenPrice          float64
	LastPrice          float64
	HighPrice          float64
	LowPrice           float64
	QuoteVolume        float64
	PriceChangePercent float64
	Count              uint64
}

var PriceChangeSignalRecord = make(map[string]*PriceChangeSignal)

func GetCoinPriceChangeSignalRecord(symbol string) *PriceChangeSignal {
	return PriceChangeSignalRecord[symbol]
}

func SetCoinPriceChangeSignalRecord(priceChangeSignal *PriceChangeSignal) {
	PriceChangeSignalRecord[priceChangeSignal.Symbol] = priceChangeSignal
}
