package routes

import (
	"sync"
)

var (
	// globalRouterManager 全局路由管理器
	globalRouterManager *RouterManager
	once                sync.Once
)

// GetGlobalRouterManager 获取全局路由管理器（单例模式）
func GetGlobalRouterManager() *RouterManager {
	once.Do(func() {
		globalRouterManager = NewRouterManager()
	})
	return globalRouterManager
}

// AutoRegister 自动注册路由组件
// 这个函数会在各个路由组件的init()函数中被调用
func AutoRegister(router Router) {
	manager := GetGlobalRouterManager()
	manager.Register(router)
}
