package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"satellite/router"
)

func main() {
	gin.SetMode(gin.DebugMode)
	ginRouter := gin.Default()
	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := ginRouter.Group("/v1")
	v1.POST("/motd", router.MotdHandler)
	err := ginRouter.Run(":8080")
	if err != nil {
		log.Fatalln(err)
		return
	}
}
