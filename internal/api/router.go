package api

import (
	"medicalchat/internal/api/middleware"

	"github.com/gin-gonic/gin"

	// 导入路由包，确保所有路由组件的init函数被执行
	"medicalchat/internal/api/routes"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {

	engine := gin.New()

	// 使用中间件
	engine.Use(gin.Logger())                // 记录日志
	engine.Use(gin.Recovery())              // 恢复中间件，防止panic导致服务崩溃
	engine.Use(middleware.CORSMiddleware()) // 跨域中间件

	// 设置所有路由
	routes.SetupRoutes(engine)

	return engine
}
