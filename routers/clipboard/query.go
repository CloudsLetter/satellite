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
	clipBoardRes := db.ClipBoard{}
	err := database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", title).Select("title, content, syntax_highlighting, expire, snapchat").First(&clipBoardRes).Error
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

	if clipBoardRes.Snapchat {
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

	ginContext.JSON(200, clipBoardRes)
}
