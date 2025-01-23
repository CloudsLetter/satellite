package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"satellite/models/db"
	"satellite/models/respon"
	"satellite/utilities/database"
	"strconv"
	"time"
)

func CallStatisticsHandler(ginContext *gin.Context) {
	day, err := strconv.Atoi(ginContext.Query("day"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}

	destination := ginContext.Query("destination")

	if destination == "call_statistics" {
		destination = "/v1/call_statistics"
	}
	if destination == "motd" {
		destination = "/v1/motd"
	}
	var resObj []respon.CallStatistics
	today := time.Now()
	firstDay := today.AddDate(0, 0, -day)
	if firstDay.Day() == today.Day() {
		tmpObj := respon.CallStatistics{}
		tmpObj.Day = today.Day()
		tmpObj.Month = int(today.Month())
		err = database.GetDBContext().Model(&db.CallsStatistics{}).Where("day = ? and destination = ?", today.Day(), destination).Count(&tmpObj.Calls).Error
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}

		resObj = append(resObj, tmpObj)
		ginContext.JSON(http.StatusOK, resObj)
		return
	}

	for i := firstDay; today.After(i); i = i.AddDate(0, 0, 1) {
		tmpObj := respon.CallStatistics{}
		tmpObj.Day = i.Day()
		tmpObj.Month = int(i.Month())
		err = database.GetDBContext().Model(&db.CallsStatistics{}).Where("day = ? AND destination = ?", i.Day(), destination).Count(&tmpObj.Calls).Error
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
		resObj = append(resObj, tmpObj)
	}
	ginContext.JSON(http.StatusOK, resObj)
	return
}
