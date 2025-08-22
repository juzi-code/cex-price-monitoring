package conf

import (
	"cex-price-monitoring/logger"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type telegramData struct {
	BotToken      string `yaml:"botToken"`
	SpotChatID    int64  `yaml:"spotChatID"`
	FuturesChatID int64  `yaml:"futuresChatID"`
}

type Config struct {
	TelegramData telegramData      `yaml:"tg"`
	LogConfig    *logger.LogConfig `yaml:"log"`
}

var (
	confPath string
	cfg      *Config
	once     sync.Once
)

func Cfg() *Config {
	return cfg
}

func LoadConfigFile() *Config {
	once.Do(func() {
		cfg = new(Config)
		// 设置默认日志配置
		cfg.LogConfig = logger.DefaultConfig()
	})
	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 如果配置文件中没有日志配置，使用默认配置
	if cfg.LogConfig == nil {
		cfg.LogConfig = logger.DefaultConfig()
	}

	return cfg
}

func initBack() {

	rootPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	file := filepath.Join(rootPath, "conf", "app.yaml")

	confPath = file

	//直接加载配置
	LoadConfigFile()
}

func init() {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录出错: %v\n", err)
		return
	}

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取可执行文件路径出错: %v\n", err)
		return
	}

	exeDir := filepath.Dir(exePath)

	// 检查工作目录是否与可执行文件所在目录相同
	if workingDir != exeDir {
		fmt.Println("当前程序可能是通过 go run 命令运行的")
	} else {
		fmt.Println("当前程序可能是通过双击可执行文件运行的")
	}

	// 获取可执行文件所在的路径
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "无法获取可执行文件路径:", err)
		os.Exit(-1)
	}

	// 获取可执行文件所在目录
	execDir := filepath.Dir(execPath)

	// 配置文件的绝对路径
	confPath = filepath.Join(execDir, "conf", "app.yaml")
	//直接加载配置
	LoadConfigFile()
}
