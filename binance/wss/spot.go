package wss

import (
	"context"
	"math"
	"strconv"
	"time"

	"cex-price-monitoring/binance"
	"cex-price-monitoring/conf"
	"cex-price-monitoring/data"
	"cex-price-monitoring/logger"
	"cex-price-monitoring/tgbot"

	binanceconnector "github.com/binance/binance-connector-go"
)

var spotWsClient *binanceconnector.WebsocketStreamClient

func init() {
	spotWsClient = binanceconnector.NewWebsocketStreamClient(true)
}

func SubSpotKLines(ctx context.Context, symbolIntervalPair map[string]string, threshold float64) {
	logger.WithFields(logger.Fields{
		"pairs": len(symbolIntervalPair), "threshold": threshold,
	}).Info("Binance-准备订阅现货K线")

	handler := func(event *binanceconnector.WsKlineEvent) {
		open, _ := strconv.ParseFloat(event.Kline.Open, 64)
		close_, _ := strconv.ParseFloat(event.Kline.Close, 64)
		high, _ := strconv.ParseFloat(event.Kline.High, 64)
		low, _ := strconv.ParseFloat(event.Kline.Low, 64)

		amplitude := (high - low) / low
		if math.Abs(amplitude) <= threshold {
			return
		}

		klineStart := time.UnixMilli(event.Kline.StartTime)
		cfg := conf.Cfg().Monitor
		if rec := data.GetSignalRecord(event.Symbol); rec != nil &&
			rec.Time.Add(cfg.Cooldown).After(klineStart) {
			return
		}

		quoteVolume, _ := strconv.ParseFloat(event.Kline.QuoteVolume, 64)
		ticker := binance.GetSpotTicker24h(event.Symbol)
		if ticker == nil {
			return
		}

		openPrice24h, _ := strconv.ParseFloat(ticker.OpenPrice, 64)
		lastPrice24h, _ := strconv.ParseFloat(ticker.LastPrice, 64)
		highPrice24h, _ := strconv.ParseFloat(ticker.HighPrice, 64)
		lowPrice24h, _ := strconv.ParseFloat(ticker.LowPrice, 64)
		quoteVolume24h, _ := strconv.ParseFloat(ticker.QuoteVolume, 64)
		pctChange24h, _ := strconv.ParseFloat(ticker.PriceChangePercent, 64)

		if quoteVolume24h < cfg.MinQuoteVolume24h {
			return
		}

		signal := data.PriceChangeSignal{
			Type:               "现货",
			Symbol:             event.Symbol,
			CexName:            "Binance",
			Interval:           event.Kline.Interval,
			OpenPrice:          open,
			ClosePrice:         close_,
			HighPrice:          high,
			LowPrice:           low,
			QuoteVolume:        quoteVolume,
			PriceChangePercent: (close_ - open) / open,
			Amplitude:          amplitude,
			TradeNum:           event.Kline.TradeNum,
			Time:               klineStart,
			CoinTicker: &data.CoinTicker{
				OpenPrice:          openPrice24h,
				LastPrice:          lastPrice24h,
				HighPrice:          highPrice24h,
				LowPrice:           lowPrice24h,
				QuoteVolume:        quoteVolume24h,
				PriceChangePercent: pctChange24h,
				Count:              uint64(ticker.Count),
			},
		}

		logger.WithFields(logger.Fields{
			"symbol": event.Symbol, "amplitude": amplitude,
		}).Info("现货振幅触发，发送通知")

		tgbot.SendPriceChangeMessage(signal, conf.Cfg().TelegramData.SpotChatID)
		data.SetSignalRecord(&signal)
	}

	errHandler := func(err error) {
		logger.WithField("error", err).Error("现货WebSocket错误")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		logger.Info("Binance-订阅现货K线（可能重连）")
		doneCh, _, err := spotWsClient.WsCombinedKlineServe(symbolIntervalPair, handler, errHandler)
		if err != nil {
			logger.WithField("error", err).Error("现货订阅失败，5s后重试")
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		select {
		case <-doneCh:
			logger.Warn("Binance-现货连接断开，准备重连")
		case <-ctx.Done():
			return
		}
		time.Sleep(5 * time.Second)
	}
}
