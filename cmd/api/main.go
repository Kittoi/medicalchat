package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"medicalchat/pkg/utils"
)

func main() {
	// 初始化配置
	if err := utils.InitConfig(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	r.Run(":8080")
}
