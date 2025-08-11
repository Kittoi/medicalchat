package service

import (
	"medicalchat/pkg/utils"
	"sync"
)

var (
	// 全局聊天服务实例
	globalChatService *ChatService
	chatServiceOnce   sync.Once
)

// GetGlobalChatService 获取全局聊天服务实例（单例模式）
func GetGlobalChatService() *ChatService {
	chatServiceOnce.Do(func() {
		config := utils.GetConfig()
		if config != nil {
			globalChatService = NewChatService(config.Chat)
		}
	})
	return globalChatService
}
