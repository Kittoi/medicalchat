package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"medicalchat/internal/api"
	"medicalchat/internal/api/routes"
	"medicalchat/internal/service"
	initApp "medicalchat/pkg/init"
)


func main() {
	// 设置 Gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化应用
	if err := initApp.InitializeApp(); err != nil {
		fmt.Printf("应用初始化失败: %v", err)
		return
	}

	chatService := service.GetGlobalChatService()
	if chatService != nil {
		routes.AutoRegister(routes.NewChatRouter(chatService))
	}

	// 创建路由引擎并设置所有路由
	r := api.SetupRouter()

	// 启动服务器
	r.Run(":8082")
}
