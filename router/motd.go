package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"satellite/model/misc"
	"satellite/model/request"
	"satellite/model/respon"
	"satellite/utilitie"
	"strconv"
	"strings"
)

func MotdHandler(ginContext *gin.Context) {

	if str := ginContext.GetHeader("Accept"); str != "application/json" {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "allow application/json only",
		})
		return
	}
	motdObj := request.MotdRequest{}

	motdObj.Host = ginContext.Query("host")
	port, err := strconv.Atoi(ginContext.Query("port"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}
	gType, err := strconv.Atoi(ginContext.Query("type"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}
	motdObj.Port = port
	motdObj.Type = gType

	if motdObj.Host == "" || motdObj.Port == 0 {
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

	jsonObj := misc.MotdJson{}
	comObj := respon.MotdResponse{}
	if motdObj.Type == 0 {
		err := utilitie.GetMotdFromJavaEdition17(&jsonObj, motdObj.Host, motdObj.Port)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "server offline or internal server error",
			})
			return
		}
		comObj = respon.MotdResponse{
			MineCraftEdition:    jsonObj.Version.Name,
			OnlinePlayers:       jsonObj.Players.Online,
			MaxPlayers:          jsonObj.Players.Max,
			Description:         jsonObj.Description.Text,
			Favicon:             jsonObj.Favicon,
			Protocol:            strconv.Itoa(jsonObj.Version.Protocol),
			Loader:              jsonObj.Modinfo.Type,
			ModCount:            len(jsonObj.Modinfo.ModList),
			EnforcesSecureChat:  jsonObj.EnforcesSecureChat,
			PreventsChatReports: jsonObj.PreventsChatReports,
		}
		comObj.Edition = "java"
		if comObj.Description == "" {
			for _, v := range jsonObj.Description.Extra {
				comObj.Description += v.Text
			}
		}

		if jsonObj.Version.Protocol == -1 && jsonObj.Version.SupportedVersions != nil {
			comObj.Protocol = ""
			for i, v := range jsonObj.Version.SupportedVersions {
				if i == 0 {
					comObj.Protocol += strconv.Itoa(v)
					continue
				}
				comObj.Protocol += "," + strconv.Itoa(v)

			}
		}

		if jsonObj.Players.Sample != nil {
			for i, v := range jsonObj.Players.Sample {
				//if i >= 10 {
				//	break
				//}
				//if v.Id == "00000000-0000-0000-0000-000000000000" {
				//	continue
				//}
				comObj.Players = append(comObj.Players, misc.Sample{
					Name: v.Name,
					Id:   v.Id,
				})
				i++
			}
		}
	}
	if motdObj.Type == 1 {
		beData, err := utilitie.GetMotdFromBedrockEdition(motdObj.Host, motdObj.Port)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{
				"message": "server offline or internal server error",
			})
		}
		comObj.MineCraftServerType = beData[0]
		comObj.ServerName = beData[1]
		comObj.Protocol = beData[2]
		comObj.Version = beData[3]
		comObj.OnlinePlayers, _ = strconv.Atoi(beData[4])
		comObj.MaxPlayers, _ = strconv.Atoi(beData[5])
		comObj.ServerId, _ = strconv.ParseInt(beData[6], 10, 64)
		comObj.Map = beData[7]
		comObj.GameMode = beData[8]
		comObj.NintendoLimited, _ = strconv.ParseBool(beData[9])
		comObj.Ipv4Port, _ = strconv.Atoi(beData[10])
		comObj.Ipv6Port, _ = strconv.Atoi(beData[11])
		comObj.Extra = beData[12]
		comObj.Edition = "bedrock"

	}

	if motdObj.Type != 0 && motdObj.Type != 1 {
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "type not support",
		})
		return
	}

	comObj.RealPlayers = len(comObj.Players)
	ginContext.JSON(http.StatusOK, comObj)
	return
}
