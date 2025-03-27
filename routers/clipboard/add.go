package clipboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"satellite/models/db"
	"satellite/models/request"
	"satellite/utilities/database"
	"satellite/utilities/misc"
	"strconv"
	"time"
)

// TODO 待设置密码

func Add(ginContext *gin.Context) {
	if str := ginContext.GetHeader("Accept"); str != "application/json" {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "allow application/json only",
		})
		return
	}
	clipboardObj := request.Clipboard{}

	if err := ginContext.ShouldBindJSON(&clipboardObj); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing Parameter",
		})
		return
	}
	timestampMill, err := strconv.ParseInt(clipboardObj.Expire, 10, 64)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter Error",
		})
		return
	}

	if matches := misc.GetACAutomaton().Search(clipboardObj.Title); len(matches) > 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "Sensitive Words Detected in Title",
		})
		return
	}

	if matches := misc.GetACAutomaton().Search(clipboardObj.Content); len(matches) > 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "Sensitive Words Detected in Content",
		})
		return
	}

	clipboardRes := db.ClipBoard{}
	err = database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", clipboardObj.Title).First(&clipboardRes).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Error",
		})
		return
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		if time.Now().After(clipboardRes.Expire) {
			if err = database.GetDBContext().Delete(&clipboardRes).Error; err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Error",
				})
				return
			}
		} else {
			ginContext.JSON(http.StatusConflict, gin.H{
				"message": "Exist",
			})
			return
		}
	}

	err = database.GetDBContext().Create(&db.ClipBoard{
		Title:    clipboardObj.Title,
		Content:  clipboardObj.Content,
		Expire:   time.UnixMilli(timestampMill),
		Snapchat: clipboardObj.Snapchat,
	}).Error
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Error",
		})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
	return
}
