package middlewares

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(func(ginContext *gin.Context) {
			ginContext.Next()
		}),
		timeout.WithResponse(func(ginContext *gin.Context) {
			ginContext.JSON(http.StatusRequestTimeout, gin.H{
				"message": "request timeout",
			})
		}),
	)
}
