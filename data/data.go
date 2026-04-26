package data

import (
	"sync"
	"time"
)

type PriceChangeSignal struct {
	Type               string
	Symbol             string
	BaseAsset          string
	CexName            string
	Interval           string
	OpenPrice          float64
	ClosePrice         float64
	HighPrice          float64
	LowPrice           float64
	QuoteVolume        float64
	PriceChangePercent float64
	Amplitude          float64
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

var (
	signalMu sync.RWMutex
	signals  = make(map[string]*PriceChangeSignal)
)

func GetSignalRecord(symbol string) *PriceChangeSignal {
	signalMu.RLock()
	defer signalMu.RUnlock()
	return signals[symbol]
}

func SetSignalRecord(s *PriceChangeSignal) {
	signalMu.Lock()
	defer signalMu.Unlock()
	signals[s.Symbol] = s
}
