package routes

import (
	"medicalchat/internal/api/handler"
	"medicalchat/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatRouter 聊天路由结构
type ChatRouter struct {
	BaseRouter
	chatHandler *handler.ChatHandler
}

// NewChatRouter 创建聊天路由
func NewChatRouter(chatService *service.ChatService) *ChatRouter {
	return &ChatRouter{
		BaseRouter:  NewBaseRouter("ChatRouter", "/api/v1"),
		chatHandler: handler.NewChatHandler(chatService),
	}
}

// RegisterRoutes 注册聊天相关路由
func (cr *ChatRouter) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group(cr.GetPrefix())
	{
		chatGroup := group.Group("/chat")
		{
			// 流式聊天
			chatGroup.POST("/stream", cr.chatHandler.ChatStream)
		}
	}
}
