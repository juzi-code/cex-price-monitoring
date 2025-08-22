package constant

import "time"

// SendMessageInterval 发送消息时间间隔
const SendMessageInterval = time.Second * (60 * 2.5)

// PriceChangeThresholdMap 价格涨幅阈值
var PriceChangeThresholdMap = map[string]float64{
	"1m": 0.04,
	"5m": 0.8,
}

// MinQuoteVolume24h 24h交易额阈值
const MinQuoteVolume24h = 100
const TimeLayoutSecond = "2006-01-02 15:04:05"
