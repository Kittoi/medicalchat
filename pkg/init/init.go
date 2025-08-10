package initApp

import (
	"fmt"

	"medicalchat/internal/repo"
	"medicalchat/pkg/utils"
)


// InitializeApp 统一处理所有初始化逻辑
func InitializeApp() error {
	// 1. 初始化配置
	if err := utils.InitConfig(); err != nil {
		return fmt.Errorf("配置初始化失败: %w", err)
	}

	// 2. 初始化日志
	if err := utils.InitLogger(); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 3. 初始化数据库连接
	if err := repo.InitDB(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	return nil
}
