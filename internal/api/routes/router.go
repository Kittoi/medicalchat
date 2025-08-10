package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Router 路由接口，所有路由组件都需要实现这个接口
type Router interface {
	// RegisterRoutes 注册路由到gin引擎
	RegisterRoutes(engine *gin.Engine)
	// GetPrefix 获取路由前缀
	GetPrefix() string
	// GetName 获取路由组件名称
	GetName() string
}

// RouterManager 路由管理器
type RouterManager struct {
	routers []Router
}

// NewRouterManager 创建新的路由管理器
func NewRouterManager() *RouterManager {
	return &RouterManager{
		routers: make([]Router, 0),
	}
}

// Register 注册路由组件
func (rm *RouterManager) Register(router Router) {
	rm.routers = append(rm.routers, router)
	log.Info().
		Str("name", router.GetName()).
		Str("prefix", router.GetPrefix()).
		Msg("路由组件注册成功")
}

// RegisterAll 注册所有路由组件到gin引擎
func (rm *RouterManager) RegisterAll(engine *gin.Engine) {
	log.Info().Msg("开始注册所有路由组件")

	for _, router := range rm.routers {
		router.RegisterRoutes(engine)
		log.Info().
			Str("name", router.GetName()).
			Str("prefix", router.GetPrefix()).
			Msg("路由注册完成")
	}

	log.Info().
		Int("count", len(rm.routers)).
		Msg("所有路由组件注册完成")
}

// GetRouters 获取所有已注册的路由组件
func (rm *RouterManager) GetRouters() []Router {
	return rm.routers
}
