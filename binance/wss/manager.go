package wss

import (
	"context"
	"sync"
	"time"

	"cex-price-monitoring/binance"
	"cex-price-monitoring/conf"
	"cex-price-monitoring/logger"
)

type MarketType string

const (
	Spot    MarketType = "spot"
	Futures MarketType = "futures"
)

type Manager struct {
	marketType MarketType

	mu        sync.RWMutex
	symbolMap map[string]string // symbol -> baseAsset

	cancel context.CancelFunc
}

func NewManager(marketType MarketType) *Manager {
	return &Manager{
		marketType: marketType,
		symbolMap:  make(map[string]string),
	}
}

func (m *Manager) Run(ctx context.Context) {
	m.refresh()
	m.startSubscriptions(ctx)

	ticker := time.NewTicker(conf.Cfg().Monitor.SymbolRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fresh := m.fetchSymbols()
			if m.hasNewSymbols(fresh) {
				logger.WithField("market", m.marketType).Info("检测到新交易对，重启订阅")
				m.updateSymbolMap(fresh)
				m.cancel()
				m.startSubscriptions(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) startSubscriptions(parentCtx context.Context) {
	subCtx, cancel := context.WithCancel(parentCtx)
	m.cancel = cancel

	cfg := conf.Cfg().Monitor
	for _, rule := range cfg.Intervals {
		pairs := m.symbolIntervalPairs(rule.Name)
		r := rule
		go func() {
			switch m.marketType {
			case Spot:
				SubSpotKLines(subCtx, pairs, r.Threshold)
			case Futures:
				SubFuturesKLines(subCtx, pairs, r.Threshold)
			}
		}()
	}
}

func (m *Manager) symbolIntervalPairs(interval string) map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pairs := make(map[string]string, len(m.symbolMap))
	for sym := range m.symbolMap {
		pairs[sym] = interval
	}
	return pairs
}

func (m *Manager) refresh() {
	fresh := m.fetchSymbols()
	m.mu.Lock()
	m.symbolMap = fresh
	m.mu.Unlock()
}

func (m *Manager) fetchSymbols() map[string]string {
	var syms []binance.Symbol
	switch m.marketType {
	case Spot:
		syms = binance.GetSpotUSDTSymbols()
	case Futures:
		syms = binance.GetFuturesUSDTSymbols()
	}
	result := make(map[string]string, len(syms))
	for _, s := range syms {
		result[s.Name] = s.BaseAsset
	}
	return result
}

func (m *Manager) hasNewSymbols(fresh map[string]string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for sym := range fresh {
		if _, exists := m.symbolMap[sym]; !exists {
			return true
		}
	}
	return false
}

func (m *Manager) updateSymbolMap(fresh map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.symbolMap = fresh
}

func (m *Manager) BaseAsset(symbol string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.symbolMap[symbol]
}
