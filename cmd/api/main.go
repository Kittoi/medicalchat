package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"medicalchat/internal/api"
	initApp "medicalchat/pkg/init"
)

func init() {
	// 初始化应用
	if err := initApp.InitializeApp(); err != nil {
		fmt.Printf("应用初始化失败: %v", err)
		return
	}
}

func main() {
	// 设置 Gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由引擎并设置所有路由
	r := api.SetupRouter()

	// 启动服务器
	r.Run(":8080")
}
