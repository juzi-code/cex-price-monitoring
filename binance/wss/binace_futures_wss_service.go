package binance

import (
	"cex-price-monitoring/binance"
	"cex-price-monitoring/conf"
	"cex-price-monitoring/constant"
	"cex-price-monitoring/data"
	"cex-price-monitoring/logger"
	"cex-price-monitoring/tgbot"
	adshaoBinance "github.com/adshao/go-binance/v2/futures"
	binanceconnector "github.com/binance/binance-connector-go"
	"math"
	"os"
	"strconv"
	"time"
)

func init() {
	httpProxy := os.Getenv("http_proxy")
	if httpProxy != "" {
		adshaoBinance.SetWsProxyUrl(httpProxy)
	}
}
func SubFuturesKLines(symbolIntervalPair map[string]string) {
	logger.WithField("pairs_count", len(symbolIntervalPair)).Info("Binance-准备订阅期货K线数据")
	wsKlineHandler := func(event *adshaoBinance.WsKlineEvent) {
		openPrice, _ := strconv.ParseFloat(event.Kline.Open, 64)
		closePrice, _ := strconv.ParseFloat(event.Kline.Close, 64)
		HighPrice, _ := strconv.ParseFloat(event.Kline.High, 64)
		LowPrice, _ := strconv.ParseFloat(event.Kline.Low, 64)
		//计算价格涨幅比例
		priceChangePercent := (closePrice - openPrice) / openPrice
		//计算价格振幅
		amplitude := (HighPrice - LowPrice) / LowPrice
		if math.Abs(amplitude) > constant.PriceChangeThresholdMap[event.Kline.Interval] {
			logger.WithField("event: %v", binanceconnector.PrettyPrint(event)).Debug("期货价格振幅触发阈值")
			coinPriceChangeSignalRecord := data.GetCoinPriceChangeSignalRecord(event.Symbol)
			klineStartTime := time.UnixMilli(event.Kline.StartTime)
			if coinPriceChangeSignalRecord != nil &&
				(*coinPriceChangeSignalRecord).Time.Add(constant.SendMessageInterval).After(klineStartTime) {
				return
			}
			highPrice, _ := strconv.ParseFloat(event.Kline.High, 64)
			lowPrice, _ := strconv.ParseFloat(event.Kline.Low, 64)
			quoteVolume, _ := strconv.ParseFloat(event.Kline.QuoteVolume, 64)
			tradeNum := event.Kline.TradeNum

			//获取日交易Ticker
			coinTicker24h := binance.GetTicker24hFutures(event.Symbol)
			if coinTicker24h != nil {
				logger.WithField("coinTicker24h", binanceconnector.PrettyPrint(coinTicker24h)).Debug("获取期货24h行情成功")
			}
			openPrice24h, _ := strconv.ParseFloat(coinTicker24h.OpenPrice, 64)
			lastPrice24h, _ := strconv.ParseFloat(coinTicker24h.LastPrice, 64)
			highPrice24h, _ := strconv.ParseFloat(coinTicker24h.HighPrice, 64)
			lowPrice24h, _ := strconv.ParseFloat(coinTicker24h.LowPrice, 64)
			quoteVolume24h, _ := strconv.ParseFloat(coinTicker24h.QuoteVolume, 64)
			priceChangePercent24h, _ := strconv.ParseFloat(coinTicker24h.PriceChangePercent, 64)
			count24H := coinTicker24h.Count

			priceChangeSignal := data.PriceChangeSignal{
				Type:               "合约",
				Symbol:             event.Symbol,
				CexName:            "Binance",
				Interval:           event.Kline.Interval,
				OpenPrice:          openPrice,
				ClosePrice:         closePrice,
				HighPrice:          highPrice,
				LowPrice:           lowPrice,
				QuoteVolume:        quoteVolume,
				PriceChangePercent: priceChangePercent,
				Amplitude:          amplitude,
				TradeNum:           tradeNum,
				Time:               klineStartTime,
				CoinTicker: &data.CoinTicker{
					Symbol:             event.Symbol,
					OpenPrice:          openPrice24h,
					LastPrice:          lastPrice24h,
					HighPrice:          highPrice24h,
					LowPrice:           lowPrice24h,
					QuoteVolume:        quoteVolume24h,
					PriceChangePercent: priceChangePercent24h,
					Count:              uint64(count24H),
				},
			}
			if quoteVolume24h >= constant.MinQuoteVolume24h {
				logger.WithFields(logger.Fields{
					"symbol":           event.Symbol,
					"amplitude":        amplitude,
					"quote_volume_24h": quoteVolume24h,
				}).Info("期货价格变动信号触发，发送通知")
				tgbot.SendPriceChangeMessage(priceChangeSignal, conf.Cfg().TelegramData.FuturesChatID)
			} else {
				logger.WithFields(logger.Fields{
					"symbol":           event.Symbol,
					"quote_volume_24h": quoteVolume24h,
					"min_required":     constant.MinQuoteVolume24h,
				}).Debug("期货交易量不足，跳过通知")
			}
			data.SetCoinPriceChangeSignalRecord(&priceChangeSignal)

		}

	}
	errHandler := func(err error) {
		logger.WithField("error", err).Error("期货WebSocket连接错误")
	}
	for {
		logger.Info("Binance-开始订阅期货K线数据（可能重连）")

		doneCh, stopCh, err := adshaoBinance.WsCombinedKlineServe(symbolIntervalPair, wsKlineHandler, errHandler)
		if err != nil {
			logger.WithField("error", err).Error("期货订阅失败，5秒后重试")
			time.Sleep(5 * time.Second)
			continue
		}

		stopCh = stopCh // 避免未使用变量警告
		logger.Info("Binance-期货K线数据订阅成功")

		<-doneCh // 等待连接关闭或出错

		logger.Warn("Binance-期货连接断开，准备重连")
		time.Sleep(5 * time.Second) // 可选：等待一段时间再重连
	}
}
