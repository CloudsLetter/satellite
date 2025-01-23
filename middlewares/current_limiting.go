package middlewares

import (
	"github.com/gin-gonic/gin"
	"satellite/models/db"
	"satellite/utilities/database"
	"sync"
	"sync/atomic"
	"time"
)

var dataPool = make(chan db.CallsStatistics, 1024)
var currentLimiting = sync.Map{}
var ticker = time.NewTicker(60 * time.Second)
var counts = int32(0)
var maxCount = int32(600)

func CurrentLimitingMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if counts < maxCount {
			atomic.AddInt32(&counts, 1)
		} else {
			ginContext.JSON(429, gin.H{
				"message": "too many requests in a minute for global",
			})
			ginContext.Abort()
			return
		}
		go func(ctx *gin.Context) {
			currentTime := time.Now()
			dataPool <- db.CallsStatistics{
				Source:      ctx.ClientIP(),
				Destination: ctx.Request.URL.Path,
				UserAgent:   ctx.GetHeader("User-Agent"),
				Minute:      currentTime.Minute(),
				Hour:        currentTime.Hour(),
				Day:         currentTime.Day(),
				Month:       int(currentTime.Month()),
				Year:        currentTime.Year(),
			}
		}(ginContext.Copy())

		key := ginContext.ClientIP() + ginContext.Request.URL.Path
		count, _ := currentLimiting.LoadOrStore(key, 0)
		if count.(int) > 10 {
			ginContext.JSON(429, gin.H{
				"message": "too many requests in a minute for local",
			})
			ginContext.Abort()
			return
		}
		currentLimiting.Store(key, count.(int)+1)
		ginContext.Next()
		atomic.AddInt32(&counts, -1)
	}
}

func WriteDataToDB() {
	for {
		select {
		case data, ok := <-dataPool:
			if !ok {
				panic(ok)
			}
			database.GetDBContext().Create(&data)
		}
	}
}
func CheckCurrentLimiting() {
	for range ticker.C {
		atomic.StoreInt32(&counts, 0)
		currentLimiting.Range(func(key, value interface{}) bool {
			currentLimiting.Delete(key)
			return true
		})
	}
}
