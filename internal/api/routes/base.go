package routes

// BaseRouter 基础路由结构
type BaseRouter struct {
	name   string
	prefix string
}

// NewBaseRouter 创建基础路由
func NewBaseRouter(name, prefix string) BaseRouter {
	return BaseRouter{
		name:   name,
		prefix: prefix,
	}
}

// GetName 获取路由组件名称
func (br *BaseRouter) GetName() string {
	return br.name
}

// GetPrefix 获取路由前缀
func (br *BaseRouter) GetPrefix() string {
	return br.prefix
}
