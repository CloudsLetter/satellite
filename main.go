package main

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"satellite/middlewares"
	"satellite/routers"
	"satellite/utilities"
)

func main() {
	utilities.InitFunc()

	gin.SetMode(gin.DebugMode)
	ginRouter := gin.Default()

	ginRouter.Use(middlewares.TimeoutMiddleware())
	ginRouter.Use(cors.Default())
	ginRouter.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
	}))
	ginRouter.Use(middlewares.CurrentLimitingMiddleware())

	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := ginRouter.Group("/v1")
	v1.GET("/motd", routers.MotdHandler)
	v1.GET("/call_statistics", routers.CallStatisticsHandler)
	err := ginRouter.Run(":8080")
	if err != nil {
		log.Fatalln(err)
		return
	}
}
