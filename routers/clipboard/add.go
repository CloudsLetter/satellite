package clipboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"satellite/macro"
	"satellite/models/db"
	"satellite/models/request"
	"satellite/utilities/database"
	"satellite/utilities/misc"
	"time"
)

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

	matches := misc.GetACAutomaton().Search(clipboardObj.Title)
	delete(matches, "")
	if len(matches) > 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status_code": macro.SensitiveWordsDetected,
			"message":     "Sensitive Words Detected in Title",
		})
		return
	}

	matches = misc.GetACAutomaton().Search(clipboardObj.Content)
	delete(matches, "")
	if len(matches) > 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status_code": macro.SensitiveWordsDetected,
			"message":     "Sensitive Words Detected in Content",
		})
		return
	}

	clipboardRes := db.ClipBoard{}

	err := database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", clipboardObj.Title).First(&clipboardRes).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Error",
		})
		return
	}

	if err == nil {
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
		Title:              clipboardObj.Title,
		Content:            clipboardObj.Content,
		Expire:             time.Now().AddDate(0, 0, clipboardObj.Expire),
		Snapchat:           clipboardObj.Snapchat,
		Password:           clipboardObj.Password,
		SyntaxHighlighting: clipboardObj.SyntaxHighlighting,
		Read:               0,
	}).Error
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Error",
		})
		return
	}
	ginContext.JSON(http.StatusCreated, gin.H{
		"message": "Success",
	})
	return
}
