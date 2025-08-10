package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"medicalchat/pkg/init"
)

func init(){
	// 初始化应用
	if err := initApp.InitializeApp(); err != nil {
		fmt.Printf("应用初始化失败: %v", err)
		return
	}
}

func main() {
	// 设置 Gin 的运行模式
	gin.SetMode(gin.ReleaseMode)

	// 创建新的 Gin 实例
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	r.Run(":8080")
}
