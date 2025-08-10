# 路由组件自动注册系统

这是一个为 Go Gin 框架设计的路由组件自动注册系统，支持模块化路由管理和自动注册功能。

## 核心组件

### 1. Router 接口
所有路由组件都需要实现 `Router` 接口：

```go
type Router interface {
    RegisterRoutes(engine *gin.Engine)  // 注册路由到gin引擎
    GetPrefix() string                  // 获取路由前缀
    GetName() string                    // 获取路由组件名称
}
```

### 2. BaseRouter 基础结构
提供路由组件的基础实现，包含通用的名称和前缀字段。

### 3. RouterManager 路由管理器
管理所有注册的路由组件，提供统一的注册和管理功能。

### 4. 自动注册机制
通过包初始化时的 `init()` 函数自动注册路由组件。

## 使用方法

### 1. 创建处理器（Handler）
```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
    return &ExampleHandler{}
}

func (h *ExampleHandler) GetExample(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "示例接口",
    })
}
```

### 2. 创建路由组件
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "medicalchat/internal/api/handler"
)

type ExampleRouter struct {
    BaseRouter
    handler *handler.ExampleHandler
}

func NewExampleRouter() *ExampleRouter {
    return &ExampleRouter{
        BaseRouter: NewBaseRouter("example", "/api/v1"),
        handler:    handler.NewExampleHandler(),
    }
}

func (er *ExampleRouter) RegisterRoutes(engine *gin.Engine) {
    group := engine.Group(er.GetPrefix())
    {
        group.GET("/example", er.handler.GetExample)
    }
}

// 重要：在 init 函数中自动注册路由
func init() {
    AutoRegister(NewExampleRouter())
}
```

### 3. 路由组件自动注册
当包被导入时，`init()` 函数会自动执行，调用 `AutoRegister()` 将路由组件注册到全局路由管理器中。

### 4. 在主程序中使用
```go
package main

import (
    "medicalchat/internal/api"
)

func main() {
    // 创建路由引擎，所有注册的路由组件会自动加载
    r := api.SetupRouter()
    r.Run(":8080")
}
```

## 已实现的路由组件

### 1. HealthRouter (/api/v1)
- `GET /api/v1/health` - 健康检查
- `GET /api/v1/version` - 获取版本信息

### 2. UserRouter (/api/v1)
- `GET /api/v1/users/` - 获取用户列表
- `POST /api/v1/users/` - 创建用户
- `GET /api/v1/users/:id` - 根据ID获取用户
- `PUT /api/v1/users/:id` - 更新用户
- `DELETE /api/v1/users/:id` - 删除用户

### 3. ChatRouter (/api/v1)
- `POST /api/v1/chat/send` - 发送消息
- `GET /api/v1/chat/history` - 获取聊天历史
- `GET /api/v1/chat/rooms/` - 获取聊天室列表
- `POST /api/v1/chat/rooms/` - 创建聊天室

## 调试功能

访问 `/router-info` 端点可以查看所有已注册的路由组件信息：

```bash
curl http://localhost:8080/router-info
```

返回示例：
```json
{
  "message": "已注册的路由组件信息",
  "routers": [
    {
      "name": "health",
      "prefix": "/api/v1"
    },
    {
      "name": "user", 
      "prefix": "/api/v1"
    },
    {
      "name": "chat",
      "prefix": "/api/v1"
    }
  ]
}
```

## 优势

1. **自动注册**：新增路由组件只需实现接口并在 init() 中注册，无需修改主程序
2. **模块化**：每个路由组件独立管理，便于维护
3. **统一管理**：所有路由组件通过RouterManager统一管理
4. **易于调试**：提供路由信息查看接口
5. **扩展性强**：新增路由组件只需按照规范实现即可

## 添加新路由组件的步骤

1. 在 `internal/api/handler/` 目录下创建处理器
2. 在 `internal/api/routes/` 目录下创建路由组件
3. 实现 `Router` 接口
4. 在 `init()` 函数中调用 `AutoRegister()`
5. 重启应用，新路由自动生效

这样设计的路由系统具有良好的可扩展性和维护性，支持快速添加新的路由模块。
