package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"satellite/model/motd"
	"satellite/model/request"
	"satellite/model/respon"
	"satellite/utilitie"
	"strings"
)

func MotdHandler(ginContext *gin.Context) {

	motdObj := request.MotdRequest{}
	err := ginContext.BindJSON(&motdObj)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "missing required fields",
		})
		return
	}
	if strings.ContainsRune(motdObj.Host, ':') {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}

	javaJsonObj := motd.MotdJavaJson{}
	if motdObj.Type == 0 {
		err = utilitie.GetMotdFromJavaVersion17(&javaJsonObj, motdObj.Host, motdObj.Port)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}
	}
	if motdObj.Type == 1 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid type",
		})
		return
	}

	if motdObj.Type == 0 {
		ginContext.JSON(http.StatusOK, respon.MotdResponse{
			Name:               javaJsonObj.Version.Name,
			OnlinePlayers:      javaJsonObj.Players.Online,
			MaxPlayers:         javaJsonObj.Players.Max,
			Description:        javaJsonObj.Description.Text,
			Favicon:            javaJsonObj.Favicon,
			Protocol:           javaJsonObj.Version.Protocol,
			Loader:             javaJsonObj.Modinfo.Type,
			ModCount:           len(javaJsonObj.Modinfo.ModList),
			EnforcesSecureChat: javaJsonObj.EnforcesSecureChat,
		})
	}

	return
}
