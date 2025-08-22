package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string `yaml:"level"`      // 日志级别: debug, info, warn, error
	Format     string `yaml:"format"`     // 日志格式: json, text
	Output     string `yaml:"output"`     // 输出方式: console, file, both
	FilePath   string `yaml:"filePath"`   // 日志文件路径
	MaxSize    int    `yaml:"maxSize"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"maxBackups"` // 保留的旧日志文件数量
	MaxAge     int    `yaml:"maxAge"`     // 日志文件保留天数
	Compress   bool   `yaml:"compress"`   // 是否压缩旧日志文件
}

// DefaultConfig 默认日志配置
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		Format:     "text",
		Output:     "both",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}
}

// Init 初始化日志系统
func Init(config *LogConfig) error {
	if config == nil {
		config = DefaultConfig()
	}

	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("无效的日志级别: %v", err)
	}
	Logger.SetLevel(level)

	// 设置日志格式
	switch config.Format {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 设置输出
	switch config.Output {
	case "console":
		Logger.SetOutput(os.Stdout)
	case "file":
		if err := setupFileOutput(config); err != nil {
			return err
		}
	case "both":
		if err := setupBothOutput(config); err != nil {
			return err
		}
	default:
		Logger.SetOutput(os.Stdout)
	}

	return nil
}

// setupFileOutput 设置文件输出
func setupFileOutput(config *LogConfig) error {
	// 确保日志目录存在
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 为文件输出设置无颜色的格式化器
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true, // 禁用颜色，避免文件中出现乱码
	})

	Logger.SetOutput(file)
	return nil
}

// setupBothOutput 设置同时输出到控制台和文件
func setupBothOutput(config *LogConfig) error {
	// 确保日志目录存在
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 创建一个自定义的Hook来处理文件输出
	Logger.AddHook(&FileHook{
		file:   file,
		config: config,
	})

	// 控制台输出保持彩色
	Logger.SetOutput(os.Stdout)
	return nil
}

// FileHook 自定义Hook，用于无颜色的文件输出
type FileHook struct {
	file   *os.File
	config *LogConfig
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	// 创建无颜色的格式化器
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true, // 禁用颜色
	}

	// 格式化日志条目
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	// 写入文件
	_, err = hook.file.Write(line)
	return err
}

// 便捷的日志方法
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Fields 类型别名，方便使用
type Fields = logrus.Fields

// WithFields 创建带字段的日志条目
func WithFields(fields Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// WithField 创建带单个字段的日志条目
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// GetCurrentTime 获取当前时间字符串
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
