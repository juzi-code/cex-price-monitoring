package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cex-price-monitoring/logger"
	"gopkg.in/yaml.v3"
)

type TelegramConfig struct {
	BotToken      string `yaml:"botToken"`
	SpotChatID    int64  `yaml:"spotChatID"`
	FuturesChatID int64  `yaml:"futuresChatID"`
}

type IntervalRule struct {
	Name      string  `yaml:"name"`
	Threshold float64 `yaml:"threshold"`
}

type MonitorConfig struct {
	SymbolRefreshInterval time.Duration  `yaml:"symbolRefreshInterval"`
	MinQuoteVolume24h     float64        `yaml:"minQuoteVolume24h"`
	Cooldown              time.Duration  `yaml:"cooldown"`
	Intervals             []IntervalRule `yaml:"intervals"`
}

type Config struct {
	TelegramData TelegramConfig    `yaml:"tg"`
	Monitor      MonitorConfig     `yaml:"monitor"`
	LogConfig    *logger.LogConfig `yaml:"log"`
}

var (
	cfg  *Config
	once sync.Once
)

func Cfg() *Config {
	return cfg
}

func init() {
	once.Do(func() {
		cfg = defaultConfig()
		loadFiles()
	})
}

func defaultConfig() *Config {
	return &Config{
		LogConfig: logger.DefaultConfig(),
		Monitor: MonitorConfig{
			SymbolRefreshInterval: 30 * time.Second,
			MinQuoteVolume24h:     1_000_000,
			Cooldown:              150 * time.Second,
			Intervals: []IntervalRule{
				{Name: "1m", Threshold: 0.04},
				{Name: "5m", Threshold: 0.008},
			},
		},
	}
}

func loadFiles() {
	execDir := executableDir()
	env := os.Getenv("APP_ENV")

	// base config
	loadYAML(filepath.Join(execDir, "conf", "app.yaml"), cfg)

	// env-specific overlay
	if env != "" {
		envFile := filepath.Join(execDir, "conf", fmt.Sprintf("app.%s.yaml", env))
		loadYAML(envFile, cfg)
	}

	if cfg.LogConfig == nil {
		cfg.LogConfig = logger.DefaultConfig()
	}
}

func loadYAML(path string, out interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("读取配置文件失败 %s: %v\n", path, err)
		}
		return
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		fmt.Printf("解析配置文件失败 %s: %v\n", path, err)
	}
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	return filepath.Dir(exe)
}
