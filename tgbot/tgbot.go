package tgbot

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cex-price-monitoring/constant"
	"cex-price-monitoring/data"
	"cex-price-monitoring/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func Init(token string) error {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	var err error
	bot, err = tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, client)
	if err != nil {
		return fmt.Errorf("初始化TelegramBot失败: %w", err)
	}
	return nil
}

func SendPriceChangeMessage(s data.PriceChangeSignal, chatID int64) {
	if bot == nil {
		logger.Error("TelegramBot未初始化")
		return
	}

	emoji := "🟢🟢🟢"
	if s.PriceChangePercent < 0 {
		emoji = "🔴🔴🔴"
	}

	coinDisplay := s.BaseAsset
	if coinDisplay == "" {
		coinDisplay = strings.TrimSuffix(s.Symbol, "USDT")
	}

	t := s.CoinTicker
	text := fmt.Sprintf(
		emoji+"***%s: 价格波动***\n"+
			"- 代币: `%s`\n"+
			"- 振幅: %s\n"+
			"- 最新价: %s\n"+
			"- 涨幅: %s(%s)\n"+
			"- 区间: ***%s***\n"+
			"- 开盘价: %s(%s)\n"+
			"- 最高价: %s(%s)\n"+
			"- 最低价: %s(%s)\n"+
			"- 交易额: %s(%s)\n"+
			"- 交易笔数: %s(%s)\n"+
			"- 时间: %s\n",
		s.Type,
		coinDisplay,
		fmtPct(s.Amplitude*100),
		fmtFloat(t.LastPrice, 8),
		fmtPct(s.PriceChangePercent*100), fmtFloat(t.PriceChangePercent, 2)+"%",
		s.CexName+"-"+s.Interval,
		fmtFloat(s.OpenPrice, 8), fmtFloat(t.OpenPrice, 8),
		fmtFloat(s.HighPrice, 8), fmtFloat(t.HighPrice, 8),
		fmtFloat(s.LowPrice, 8), fmtFloat(t.LowPrice, 8),
		fmtAmount(s.QuoteVolume), fmtAmount(t.QuoteVolume),
		fmtAmount(float64(s.TradeNum)), fmtAmount(float64(t.Count)),
		fmtTime(s.Time),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📈现货交易",
				fmt.Sprintf("https://www.binance.com/zh-CN/trade/%s?type=spot", s.Symbol)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("⚔️合约交易",
				fmt.Sprintf("https://www.binance.com/zh-CN/futures/%s", s.Symbol)),
		),
	)

	resp, err := bot.Send(msg)
	if err != nil {
		logger.WithFields(logger.Fields{"chat_id": chatID, "symbol": s.Symbol, "error": err}).
			Error("发送Telegram消息失败")
		return
	}
	logger.WithFields(logger.Fields{
		"chat_id": resp.Chat.ID, "message_id": resp.MessageID, "symbol": s.Symbol,
	}).Info("Telegram消息发送成功")
}

func fmtPct(f float64) string {
	return fmtFloat(f, 2) + "%"
}

func fmtFloat(f float64, prec int) string {
	s := strconv.FormatFloat(f, 'f', prec, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func fmtAmount(v float64) string {
	switch {
	case v >= 1e9:
		return fmt.Sprintf("%.2fB", v/1e9)
	case v >= 1e6:
		return fmt.Sprintf("%.2fM", v/1e6)
	case v >= 1e3:
		return fmt.Sprintf("%.2fK", v/1e3)
	default:
		return fmt.Sprintf("%.2f", v)
	}
}

func fmtTime(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return t.Format(constant.TimeLayoutSecond)
	}
	return t.In(loc).Format(constant.TimeLayoutSecond)
}
