package tgbot

import (
	"cex-price-monitoring/conf"
	"cex-price-monitoring/constant"
	"cex-price-monitoring/data"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewTgBot() *tgbotapi.BotAPI {

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{
		Transport: transport,
	}

	// 使用自定义的 http.Client 初始化 bot
	bot, err := tgbotapi.NewBotAPIWithClient(conf.Cfg().TelegramData.BotToken, tgbotapi.APIEndpoint, client)
	if err != nil {
		log.Panicf("无法初始化BotAPI: %v", err)
	}

	bot.Debug = true
	return bot
}

func SendPriceChangeMessage(priceChangeSignal data.PriceChangeSignal, telegramChatID int64) {
	_type := priceChangeSignal.Type
	symbol := priceChangeSignal.Symbol
	cexName := priceChangeSignal.CexName
	interval := priceChangeSignal.Interval
	openPrice := priceChangeSignal.OpenPrice
	//closePrice := priceChangeSignal.ClosePrice
	highPrice := priceChangeSignal.HighPrice
	lowPrice := priceChangeSignal.LowPrice
	quoteVolume := priceChangeSignal.QuoteVolume
	PriceChangePercent := priceChangeSignal.PriceChangePercent
	TradeNum := priceChangeSignal.TradeNum
	Time := priceChangeSignal.Time

	OpenPrice24h := priceChangeSignal.CoinTicker.OpenPrice
	LastPrice24h := priceChangeSignal.CoinTicker.LastPrice
	HighPrice24h := priceChangeSignal.CoinTicker.HighPrice
	LowPrice24h := priceChangeSignal.CoinTicker.LowPrice
	quoteVolume24h := priceChangeSignal.CoinTicker.QuoteVolume
	PriceChangePercent24h := priceChangeSignal.CoinTicker.PriceChangePercent
	Count24h := priceChangeSignal.CoinTicker.Count

	// 使用 fmt.Sprintf 格式化消息字符串
	emojiStr := "🟢🟢🟢"
	if PriceChangePercent < 0 {
		emojiStr = "🔴🔴🔴"
	}
	messageStr := fmt.Sprintf(
		emojiStr+
			"***价格波动***\n"+
			"- 类型: %s\n"+
			"- 代币: `%s`\n"+
			"- 区间: ***%s***\n"+
			"- 涨幅: %s(%s)\n"+
			"- 最新价: %s\n"+
			"- 开盘价: %s(%s)\n"+
			"- 最高价: %s(%s)\n"+
			"- 最低价: %s(%s)\n"+
			"- 交易额: %s(%s)\n"+
			"- 交易笔数: %s(%s)\n"+
			"- 时间: %s\n",
		_type,
		strings.Replace(symbol, "USDT", "", -1),
		cexName+"-"+interval,
		formatFloatToStr(PriceChangePercent*100, 2)+"%",
		formatFloatToStr(PriceChangePercent24h, 2)+"%",
		formatFloatToStr(LastPrice24h, 8),
		formatFloatToStr(openPrice, 8),
		formatFloatToStr(OpenPrice24h, 8),
		formatFloatToStr(highPrice, 8),
		formatFloatToStr(HighPrice24h, 8),
		formatFloatToStr(lowPrice, 8),
		formatFloatToStr(LowPrice24h, 8),
		formatAmount(quoteVolume), formatAmount(quoteVolume24h),
		formatAmount(float64(TradeNum)), formatAmount(float64(Count24h)),
		formatDateStr(Time),
	)
	msg := tgbotapi.NewMessage(telegramChatID, messageStr)
	msg.ParseMode = tgbotapi.ModeMarkdown
	message, _err := NewTgBot().Send(msg)
	if _err != nil {
		log.Panicf("发送消息失败: %v", _err)
	} else {
		log.Printf("Message sent to chatID %d: %s", message.Chat.ID, message.Text)
	}
}

func formatAmount(amount float64) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("%.2fB", amount/1000000000.0) // 转换为 Billion
	} else if amount >= 1000000 {
		return fmt.Sprintf("%.2fM", amount/1000000.0) // 转换为 Million
	} else if amount >= 1000 {
		return fmt.Sprintf("%.2fK", amount/1000.0) // 转换为 Thousand
	}
	return fmt.Sprintf("%.2f", amount) // 默认返回原始数值
}

func formatDateStr(_time time.Time) string {
	// 设置 UTC+8 时区
	loc, err := time.LoadLocation("Asia/Shanghai") // 上海时间（UTC+8）
	if err != nil {
		fmt.Println("加载时区失败:", err)
		return ""
	}

	// 转换为 UTC+8 时区的时间
	utcPlus8Time := _time.In(loc)

	// 格式化时间为字符串
	return utcPlus8Time.Format(constant.TimeLayoutSecond) // 格式化为标准的时间字符串
}

func formatFloatToStr(f float64, prec int) string {
	// 格式化浮动数字，保留8位小数
	str := strconv.FormatFloat(f, 'f', prec, 64)

	// 去掉末尾的0和小数点（如果没有有效数字）
	str = strings.TrimRight(str, "0")

	// 如果最后一个字符是小数点，去掉小数点
	str = strings.TrimRight(str, ".")
	return str
}
