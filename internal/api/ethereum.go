package api

import (
	"medicalchat/internal/api/handler"
	"medicalchat/internal/api/routes"
	"medicalchat/internal/service"
	"medicalchat/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// SetupEthereumRouter 设置以太坊路由
func SetupEthereumRouter() *gin.Engine {
	engine := gin.New()

	// 使用中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 获取配置
	config := utils.GetConfig()

	// 创建以太坊服务
	ethService, err := service.NewEthereumService(config.Ethereum, log.Logger)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Ethereum service, some features may not work")
		panic("Eth Client Error")
	}

	// 创建处理器并注册路由
	if ethService != nil {
		ethHandler := handler.NewEthereumHandler(ethService, log.Logger)
		routes.RegisterEthereumRoutes(engine, ethHandler)
		log.Info().Msg("Ethereum routes registered successfully")

		// 设置清理函数（这里可以考虑使用context来处理）
		// 注意：在实际应用中，应该在程序退出时调用ethService.Close()
	} else {
		log.Warn().Msg("Ethereum service not available, routes not registered")
	}

	return engine
}
