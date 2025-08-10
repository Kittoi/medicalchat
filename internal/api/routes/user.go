package routes

import (
	"medicalchat/internal/api/handler"
	"github.com/gin-gonic/gin"
)

// UserRouter 用户路由
type UserRouter struct {
	BaseRouter
	handler *handler.UserHandler
}

// NewUserRouter 创建用户路由
func NewUserRouter() *UserRouter {
	return &UserRouter{
		BaseRouter: NewBaseRouter("user", "/api/v1"),
		handler:    handler.NewUserHandler(),
	}
}

// RegisterRoutes 注册路由
func (ur *UserRouter) RegisterRoutes(engine *gin.Engine) {
	group := engine.Group(ur.GetPrefix())
	{
		users := group.Group("/users")
		{
			users.GET("/", ur.handler.GetUsers)
			users.POST("/", ur.handler.CreateUser)
			users.GET("/:id", ur.handler.GetUserByID)
			users.PUT("/:id", ur.handler.UpdateUser)
			users.DELETE("/:id", ur.handler.DeleteUser)
		}
	}
}

// init 包初始化时自动注册路由
func init() {
	AutoRegister(NewUserRouter())
}
