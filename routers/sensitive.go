package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"satellite/macro"
	"satellite/models/request"
	"satellite/utilities/misc"
)

// TODO 等待优化存储结构，内存中的结构

// Sensitive is a handler for sensitive words
// @BasePath /api/v1
// @Summary Check the sensitive words
// @Description Check the sensitive words
// @Tags sensitive
// @Accept json
// @Produce json
// @Param keywords body request.Sensitive true "The keywords of the sensitive words" Example: {"keywords
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /sensitive [post]
func Sensitive(ginContext *gin.Context) {
	if str := ginContext.GetHeader("Accept"); str != "application/json" {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status_code": macro.AcceptApplicationJsonOnly,
			"message":     "allow application/json only",
		})
		return
	}

	resObj := request.Sensitive{}

	if err := ginContext.ShouldBindJSON(&resObj); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"status_code": macro.MissingParameter,
			"message":     "Missing Parameter",
		})
		return
	}

	matches := misc.GetACAutomaton().Search(resObj.Keywords)
	delete(matches, "")
	if len(matches) > 0 {
		ginContext.JSON(http.StatusOK, gin.H{
			"status_code":     macro.SensitiveWordsDetected,
			"message":         "Sensitive Words Detected",
			"sensitive_words": matches,
		})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"status_code": macro.NoSensitiveWordsDetected,
		"message":     "No Sensitive Words Detected",
	})
	return
}
