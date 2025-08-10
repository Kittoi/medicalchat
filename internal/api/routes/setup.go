package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
// 这个函数会自动注册所有已经通过AutoRegister注册的路由组件
func SetupRoutes(engine *gin.Engine) {
	// 获取全局路由管理器
	manager := GetGlobalRouterManager()

	// 注册所有路由组件
	manager.RegisterAll(engine)
}
