package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger 初始化日志配置
func InitLogger() error {
	config := GetConfig()

	// 设置时间格式
	zerolog.TimeFieldFormat = time.RFC3339

	// 创建日志目录
	if err := os.MkdirAll(config.Log.Dir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建日志文件
	currentTime := time.Now()
	logFileName := filepath.Join(config.Log.Dir, fmt.Sprintf("%s.log", currentTime.Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %w", err)
	}

	// 同时输出到控制台和文件
	multiWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}, logFile)

	// 初始化全局日志记录器
	Logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()

	// 替换全局 logger
	log.Logger = Logger

	return nil
}