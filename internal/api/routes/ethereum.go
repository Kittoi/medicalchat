package routes

import (
	"medicalchat/internal/api/handler"

	"github.com/gin-gonic/gin"
)

// RegisterEthereumRoutes 注册以太坊相关路由
func RegisterEthereumRoutes(r *gin.Engine, ethHandler *handler.EthereumHandler) {
	// 区块链网络状态相关路由
	blockchain := r.Group("/api/blockchain")
	{
		// 获取网络状态
		blockchain.GET("/status", ethHandler.GetNetworkStatus)
	}

	// 存证相关路由
	certificate := r.Group("/api/certificate")
	{
		// 创建存证
		certificate.POST("/", ethHandler.CreateCertificate)

		// 验证存证
		certificate.POST("/verify", ethHandler.VerifyCertificate)

		// 获取存证列表 (支持状态过滤)
		certificate.GET("/", ethHandler.GetCertificates)

		// 根据ID获取存证详情
		certificate.GET("/:id", ethHandler.GetCertificateByID)
	}
}
