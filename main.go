package main

import (
	"github.com/gin-gonic/gin"
	"yinkn.cn/ws-server/controller"
	"yinkn.cn/ws-server/service"
)

func main() {
	defer service.DefaultClientManager.Close()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ws", controller.Ws)
	r.POST("/send", controller.Send)
	r.Run(":8080")
}
