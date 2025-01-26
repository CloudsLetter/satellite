package main

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"satellite/docs"
	"satellite/middlewares"
	"satellite/routers"
	"satellite/utilities"
)

func main() {
	utilities.InitFunc()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://55268fcd7867979562654ef2f4660340@o926347.ingest.us.sentry.io/4508697946423296",
		EnableTracing:    true,
		TracesSampleRate: 0.2,
	}); err != nil {
		panic(err)
	}
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
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := ginRouter.Group("/api/v1")
	v1.GET("/motd", routers.MotdHandler)
	v1.GET("/call_statistics", routers.CallStatisticsHandler)
	v1.GET("/total_calls_sse", routers.TotalCallsSSE)
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := ginRouter.Run(":8080")
	if err != nil {
		log.Fatalln(err)
		return
	}
}
