package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"

	"medicalchat/internal/models"
)

var GlobalConfig models.Config

// InitConfig 初始化配置
func InitConfig() error {
	// 获取项目根目录
	_, b, _, _ := runtime.Caller(0)
	// 获取项目根目录的绝对路径
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(b)))

	// 设置配置文件路径
	viper.AddConfigPath(filepath.Join(basePath, "configs"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将配置解析到结构体中
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// GetConfig 获取配置
func GetConfig() *models.Config {
	return &GlobalConfig
}
