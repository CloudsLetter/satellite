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

func extraRawRecursion(extra []misc.ExtraObjectOrString, descriptionRaw *[]respon.DescriptionRaw, depth int) {
	for _, v := range extra {
		result := respon.DescriptionRaw{}
		if v.ExtraString != "" {
			result.Text += v.ExtraString
		}
		if v.ExtraObject.Text != "" {
			result.Text += v.ExtraObject.Text
		}
		if v.ExtraObject.Color != "" {
			result.Color = v.ExtraObject.Color
		}
		if v.ExtraObject.Bold {
			result.Bold = v.ExtraObject.Bold
		}
		result.Depth = depth
		*descriptionRaw = append(*descriptionRaw, result)
		if v.ExtraObject.Extra != nil {
			extraRawRecursion(v.ExtraObject.Extra, descriptionRaw, depth+1)
		}
	}
}

func extraRecursion(extra []misc.ExtraObjectOrString) string {
	result := ""
	for _, v := range extra {
		if v.ExtraString != "" {
			result += v.ExtraString
		}
		if v.ExtraObject.Text != "" {
			result += v.ExtraObject.Text
		}
		if v.ExtraObject.Extra != nil {
			result += extraRecursion(v.ExtraObject.Extra)
		}
	}
	return result
}

func MotdHandler(ginContext *gin.Context) {

	if str := ginContext.GetHeader("Accept"); str != "application/json" {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "allow application/json only",
		})
		return
	}
	reqObj := request.MotdRequest{}

	reqObj.Host = ginContext.Query("host")
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
	raw, err := strconv.ParseBool(ginContext.Query("raw"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}
	desRaw, err := strconv.ParseBool(ginContext.Query("desraw"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}

	anonymous, err := strconv.ParseBool(ginContext.Query("anonymous"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}
	reqObj.Port = port
	reqObj.Type = gType
	reqObj.Raw = raw
	reqObj.DesRaw = desRaw
	reqObj.Anonymous = anonymous
	if reqObj.Host == "" || reqObj.Port == 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "missing required fields",
		})
		return
	}

	if strings.ContainsRune(reqObj.Host, ':') {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"message": "illegal input",
		})
		return
	}

	comObj := respon.MotdResponse{}
	if reqObj.Type == 0 {
		motdObj := misc.JavaMotdJson17{}
		err = utilitie.GetMotdFromJavaEdition17(&motdObj, reqObj.Host, reqObj.Port)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "server offline",
			})
			return
		}
		if reqObj.Raw {
			motdObj.Latency = 0
			ginContext.JSON(http.StatusOK, motdObj)
			return
		}
		comObj = respon.MotdResponse{
			MineCraftEdition:    motdObj.Version.Name,
			OnlinePlayers:       motdObj.Players.Online,
			MaxPlayers:          motdObj.Players.Max,
			Favicon:             motdObj.Favicon,
			Protocol:            strconv.Itoa(motdObj.Version.Protocol),
			Loader:              motdObj.Modinfo.Type,
			ModCount:            len(motdObj.Modinfo.ModList),
			EnforcesSecureChat:  motdObj.EnforcesSecureChat,
			PreventsChatReports: motdObj.PreventsChatReports,
			Latency:             motdObj.Latency,
		}
		comObj.Edition = "java"

		if motdObj.Description.DescriptionString != "" {
			comObj.Description = motdObj.Description.DescriptionString
		}
		if motdObj.Description.DescriptionObject.Text != "" {
			comObj.Description = motdObj.Description.DescriptionObject.Text
		}

		if motdObj.Description.DescriptionString == "" && motdObj.Description.DescriptionObject.Text == "" {
			if reqObj.DesRaw {
				var descriptionRaw []respon.DescriptionRaw
				extraRawRecursion(motdObj.Description.DescriptionObject.ExtraObj, &descriptionRaw, 0)
				comObj.DescriptionRaw = descriptionRaw
			} else {
				comObj.Description = extraRecursion(motdObj.Description.DescriptionObject.ExtraObj)
			}
		}

		if motdObj.Version.Protocol == -1 && motdObj.Version.SupportedVersions != nil {
			comObj.Protocol = ""
			for i, v := range motdObj.Version.SupportedVersions {
				if i == 0 {
					comObj.Protocol += strconv.Itoa(v)
					continue
				}
				comObj.Protocol += "," + strconv.Itoa(v)

			}
		}

		if motdObj.Players.Sample != nil {
			for i, v := range motdObj.Players.Sample {
				//if i >= 10 {
				//	break
				//}
				if reqObj.Anonymous {
					if v.Id == "00000000-0000-0000-0000-000000000000" {
						continue
					}
				}

				comObj.Players = append(comObj.Players, misc.Sample{
					Name: v.Name,
					Id:   v.Id,
				})
				i++
			}
		}
	}
	if reqObj.Type == 1 {
		motdObj := misc.BedrockMotd{}
		err = utilitie.GetMotdFromBedrockEdition(&motdObj, reqObj.Host, reqObj.Port)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{
				"message": "server offline",
			})
			return
		}
		if reqObj.Raw {
			motdObj.Latency = 0
			ginContext.JSON(http.StatusOK, motdObj)
			return
		}
		comObj.MineCraftServerType = motdObj.GameName
		comObj.ServerName = motdObj.HostName
		comObj.Protocol = motdObj.Protocol
		comObj.Version = motdObj.Version
		comObj.OnlinePlayers = motdObj.Players
		comObj.MaxPlayers = motdObj.MaxPlayers
		comObj.ServerId = motdObj.ServerId
		comObj.Map = motdObj.Map
		comObj.GameMode = motdObj.GameMode
		comObj.NintendoLimited = motdObj.NintendoLimited
		comObj.Ipv4Port = motdObj.IPv4Port
		comObj.Ipv6Port = motdObj.IPv6Port
		comObj.Extra = motdObj.Extra
		comObj.Edition = "bedrock"
		comObj.Latency = motdObj.Latency
	}

	if reqObj.Type == 2 || reqObj.Type == 3 {
		motdObj := misc.JavaMotd16{}
		err = utilitie.GetMotdFromJavaEdition16(&motdObj, reqObj.Host, reqObj.Port, reqObj.Type)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{
				"message": "server offline",
			})
			return
		}
		if reqObj.Raw {
			motdObj.Latency = 0
			ginContext.JSON(http.StatusOK, motdObj)
			return
		}
		comObj.Protocol = motdObj.Protocol
		comObj.Version = motdObj.Version
		comObj.Description = motdObj.MessageDaily
		comObj.OnlinePlayers = motdObj.OnlinePlayers
		comObj.MaxPlayers = motdObj.MaxPlayers
		comObj.Latency = motdObj.Latency
		comObj.Edition = "java"
	}

	if reqObj.Type == 4 {

	}

	if reqObj.Type > 4 {
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "type not support",
		})
		return
	}

	comObj.RealPlayers = len(comObj.Players)
	ginContext.JSON(http.StatusOK, comObj)
	return
}
