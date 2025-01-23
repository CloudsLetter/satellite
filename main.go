package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"satellite/middlewares"
	"satellite/routers"
	"satellite/utilities/database"
)

func main() {

	if runtime.GOOS == "windows" {
		database.InitDB("host=127.0.0.1 user=postgres password=88323095.zxc dbname=satellitetmp port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	} else {
		database.InitDB("host=postgres user=postgres password=88323095.zxc dbname=satellite port=5432 sslmode=disable TimeZone=Asia/Shanghai")
	}
	database.AutoMigrate(database.GetDBContext())
	for i := 0; i < 3; i++ {
		go middlewares.WriteDataToDB()
	}
	go middlewares.CheckCurrentLimiting()
	gin.SetMode(gin.DebugMode)
	ginRouter := gin.Default()
	ginRouter.Use(cors.Default())
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
