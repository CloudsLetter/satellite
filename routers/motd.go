package routers

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"satellite/models/minecraft"
	"satellite/models/request"
	"satellite/models/respon"
	"satellite/utilities/database"
	"satellite/utilities/misc"
	"strconv"
	"strings"
	"time"
)

// MotdHandler
// @BasePath /api/v1
// @Summary Get the message of the day from a Minecraft server
// @Description Get the message of the day from a Minecraft server
// @Tags motd
// @Accept json
// @Produce json
// @Param host query string true "The host of the Minecraft server"
// @Param port query int true "The port of the Minecraft server"
// @Param type query int true "The type of the Minecraft server"
// @Param raw query bool true "Whether to return the raw MOTD"
// @Param desraw query bool true "Whether to return the raw description"
// @Param anonymous query bool true "Whether to return anonymous players"
// @Success 200 {object} respon.MotdResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /motd [get]
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

	hub := sentrygin.GetHubFromContext(ginContext)
	if hub == nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	comObj := respon.MotdResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	val, err := database.GetRedisContext().Get(ctx, reqObj.Host+ginContext.Query("port")+ginContext.Query("type")+ginContext.Query("raw")+ginContext.Query("desraw")+ginContext.Query("anonymous")).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("CallStatistics", "redis")
				hub.CaptureException(err)
			})
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
	} else {
		if val == nil {
			ginContext.JSON(http.StatusOK, gin.H{
				"message": "empty",
			})
			return
		}

		dec := gob.NewDecoder(bytes.NewReader(val))
		if err := dec.Decode(&comObj); err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("CallStatistics", "decode")
				hub.CaptureException(err)
			})
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
		ginContext.JSON(http.StatusOK, comObj)
		return
	}

	if reqObj.Type == 0 {
		motdObj := minecraft.JavaMotdJson17{}
		err = misc.GetMotdFromJavaEdition17(&motdObj, reqObj.Host, reqObj.Port)
		if err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("motd", "java17")
				hub.CaptureException(err)
			})
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

		if motdObj.Description.DescriptionString == "" && motdObj.Description.DescriptionObject.Text == "" && motdObj.Description.DescriptionObject.ExtraObj != nil {
			if reqObj.DesRaw {
				var descriptionRaw []respon.DescriptionRaw
				misc.ExtraRawRecursion(motdObj.Description.DescriptionObject.ExtraObj, &descriptionRaw, 0)
				comObj.DescriptionRaw = descriptionRaw
			} else {
				comObj.Description = misc.ExtraRecursion(motdObj.Description.DescriptionObject.ExtraObj)
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
				comObj.Players = append(comObj.Players, minecraft.Sample{
					Name: v.Name,
					Id:   v.Id,
				})
				i++
			}
		}
	}
	if reqObj.Type == 1 {
		motdObj := minecraft.BedrockMotd{}
		err = misc.GetMotdFromBedrockEdition(&motdObj, reqObj.Host, reqObj.Port)
		if err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("motd", "bedrock")
				hub.CaptureException(err)
			})
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
		motdObj := minecraft.JavaMotd16{}
		err = misc.GetMotdFromJavaEdition16(&motdObj, reqObj.Host, reqObj.Port, reqObj.Type)
		if err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("motd", "java16")
				hub.CaptureException(err)
			})
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
		motdObj := minecraft.JavaMotd13{}
		err = misc.GetMotdFromJavaVersion13(&motdObj, reqObj.Host, reqObj.Port)
		if err != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("motd", "java13")
				hub.CaptureException(err)
			})
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
		comObj.Description = motdObj.Motd
		comObj.OnlinePlayers = motdObj.OnlinePlayers
		comObj.MaxPlayers = motdObj.MaxPlayers
		comObj.Latency = motdObj.Latency
		comObj.Edition = "java"
	}

	if reqObj.Type > 4 {
		ginContext.JSON(http.StatusOK, gin.H{
			"message": "type not support",
		})
		return
	}

	comObj.RealPlayers = len(comObj.Players)

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(comObj); err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("gob", "encode")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	if err = database.GetRedisContext().Set(ctx, reqObj.Host+ginContext.Query("port")+ginContext.Query("type")+ginContext.Query("raw")+ginContext.Query("desraw")+ginContext.Query("anonymous"), buf.Bytes(), 3*time.Minute).Err(); err != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("redis", "set")
			hub.CaptureException(err)
		})
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
		return
	}
	ginContext.JSON(http.StatusOK, comObj)
	return
}
