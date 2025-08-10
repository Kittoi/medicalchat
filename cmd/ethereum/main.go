package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"medicalchat/internal/api"
	initApp "medicalchat/pkg/init"
)

func init() {
	// 初始化以太坊应用
	if err := initApp.InitializeEthereumApp(); err != nil {
		fmt.Printf("以太坊应用初始化失败: %v", err)
		return
	}
}

func main() {
	// 设置 Gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	// 创建以太坊路由引擎
	r := api.SetupEthereumRouter()

	// 启动服务器
	log.Info().Msg("Starting Medical Chat Ethereum Toolkit server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
