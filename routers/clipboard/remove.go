package clipboard

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"satellite/models/db"
	"satellite/utilities/database"
)

func Remove(ginContext *gin.Context) {
	title := ginContext.Query("title")
	password := ginContext.Query("password")
	clipBoardRes := db.ClipBoard{}
	err := database.GetDBContext().Model(&db.ClipBoard{}).Where("title = ?", title).First(&clipBoardRes).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		ginContext.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}

	if clipBoardRes.Password == "" || clipBoardRes.Password == password {
		if err = database.GetDBContext().Delete(&clipBoardRes).Error; err != nil {
			ginContext.JSON(500, gin.H{
				"message": "Internal Error",
			})
			return
		}
		ginContext.JSON(202, gin.H{
			"message": "Delete Success",
		})
		return

	}

	ginContext.JSON(403, gin.H{
		"message": "Password Error",
	})
	return
}
