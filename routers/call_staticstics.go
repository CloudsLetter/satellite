package routers

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"satellite/models/db"
	"satellite/models/respon"
	"satellite/utilities/database"
	"time"
)

// CallStatisticsHandler is a handler for call statistics
// @BasePath /api/v1
// @Summary Get the call statistics
// @Description Get the call statistics
// @Tags call_statistics
// @Accept json
// @Produce json
// @Param destination query string true "The destination of the call" Example: /ping
// @Success 200 {object} respon.CallStatistics
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /call_statistics [get]
func CallStatisticsHandler(ginContext *gin.Context) {
	destination := ginContext.Query("destination")
	if destination == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}

	hub := sentrygin.GetHubFromContext(ginContext)
	if hub == nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	var resObj respon.CallStatistics

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	val, errs := database.GetRedisContext().Get(ctx, destination).Bytes()
	if errs != nil {
		if !errors.Is(errs, redis.Nil) {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("CallStatistics", "redis")
				hub.CaptureException(errs)
			})
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
	} else {
		if val == nil {
			ginContext.JSON(http.StatusNotFound, gin.H{
				"message": "empty",
			})
			return
		}
		dec := gob.NewDecoder(bytes.NewReader(val))
		if err := dec.Decode(&resObj); err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("CallStatistics", "decode")
				hub.CaptureException(err)
			})
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
		ginContext.JSON(http.StatusOK, resObj)
		return
	}

	var err error
	today := time.Now()
	if err = database.GetDBContext().Model(&db.CallsStatistics{}).Select("hour, minute, COUNT(*) AS count").Where("(hour = ? AND minute >= ?) OR (hour = ? AND minute <= ?)", today.Add(-time.Duration(1)*time.Hour).Hour(), today.Add(-time.Duration(1)*time.Hour).Minute(), today.Add(0).Hour(), today.Add(0).Minute()).Where("destination = ?", destination).Group("hour,minute").Order("hour,minute").Find(&resObj.MinuteStruct).Error; err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("CallStatistics", "today")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}

	if err = database.GetDBContext().Model(&db.CallsStatistics{}).Select("day, hour, COUNT(*) AS count").Where("(day = ? AND hour >= ?) OR (day = ? AND hour <= ?)", today.Add(-time.Duration(24)*time.Hour).Day(), today.Add(-time.Duration(24)*time.Hour).Hour(), today.Add(-time.Duration(0)*time.Hour).Day(), today.Add(0).Hour()).Where("destination = ?", destination).Group("day,hour").Order("day,hour").Find(&resObj.HourStruct).Error; err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("CallStatistics", "today")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}

	if err = database.GetDBContext().Model(&db.CallsStatistics{}).Select("month, day, COUNT(*) AS count").Where("(month = ? AND day >= ?) OR (month = ? AND day <= ?)", today.AddDate(0, 0, -30).Month(), today.AddDate(0, 0, -30).Day(), today.AddDate(0, 0, 0).Month(), today.AddDate(0, 0, 0).Day()).Where("destination = ?", destination).Group("month,day").Order("month,day").Find(&resObj.DayStruct).Error; err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("CallStatistics", "today")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}

	if err = database.GetDBContext().Model(&db.CallsStatistics{}).Select("year, month, COUNT(*) AS count").Where("(year = ? AND month >= ?) OR (year = ? AND month <= ?)", today.AddDate(0, 12, 0).Year(), today.AddDate(0, 12, 0).Month(), today.AddDate(0, 0, 0).Year(), today.AddDate(0, 0, 0).Month()).Where("destination = ?", destination).Group("year,month").Order("year,month").Find(&resObj.MonthStruct).Error; err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("CallStatistics", "today")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(resObj)
	if err = database.GetRedisContext().Set(ctx, destination, buf.Bytes(), 3*time.Minute).Err(); err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("CallStatistics", "redis")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	ginContext.JSON(http.StatusOK, resObj)
	return
}
