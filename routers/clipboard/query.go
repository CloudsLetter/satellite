package clipboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"satellite/models/db"
	"satellite/utilities/database"
	"time"
)

func Query(ginContext *gin.Context) {
	title := ginContext.Query("title")
	password := ginContext.Query("password")
	clipBoardRes := db.ClipBoard{}
	err := database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", title).First(&clipBoardRes).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		ginContext.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	} else if err != nil {
		ginContext.JSON(500, gin.H{
			"message": "Internal Error",
		})
		return
	}
	clipBoardRes.Read+=1
	database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", title).Update("read", clipBoardRes.Read)

	if clipBoardRes.Snapchat && clipBoardRes.Read > 1 {
		// delete this record
		if err = database.GetDBContext().Delete(&clipBoardRes).Error; err != nil {
			ginContext.JSON(500, gin.H{
				"message": "Internal Error",
			})
			return

		}
	}
	if time.Now().After(clipBoardRes.Expire) {
		// delete this record

		if err = database.GetDBContext().Delete(&clipBoardRes).Error; err != nil {
			ginContext.JSON(500, gin.H{
				"message": "Internal Error",
			})
			return

		}

		ginContext.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}
	if clipBoardRes.Password != "" && password == "" {
		ginContext.JSON(403, gin.H{
			"message": "Password Required",
		})
	}

	if clipBoardRes.Password != password {
		ginContext.JSON(403, gin.H{
			"message": "Password Error",
		})
		return
	}

	ginContext.JSON(200, clipBoardRes)
}
