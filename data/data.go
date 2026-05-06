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
	signalMu   sync.RWMutex
	signals    = make(map[string]*PriceChangeSignal)
	notifiedAt = make(map[string]time.Time)
)

// CheckAndSet atomically checks the cooldown and records the signal if not in cooldown.
// Returns true if the signal was recorded (caller should send notification).
func CheckAndSet(s *PriceChangeSignal, cooldown time.Duration) bool {
	signalMu.Lock()
	defer signalMu.Unlock()
	if t, exists := notifiedAt[s.Symbol]; exists && time.Since(t) < cooldown {
		return false
	}
	notifiedAt[s.Symbol] = time.Now()
	signals[s.Symbol] = s
	return true
}
