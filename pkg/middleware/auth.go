package middleware

import (
	"gServ/core/config"
	"gServ/pkg/gserv"
	"gServ/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func CodeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Auth-Code")
		if authorization == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		authorization_slice := strings.Split(authorization, " ")
		if len(authorization_slice) != 2 || authorization_slice[0] != "Bearer" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		auth_code := authorization_slice[1]

		if auth_code != config.GetConfig().Server.AuthCode {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
	}
}

func PlayerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		authorization_slice := strings.Split(authorization, " ")
		if len(authorization_slice) != 2 || authorization_slice[0] != "Bearer" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		token := authorization_slice[1]

		auth_player, err := playerAuth(token)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		setAuthPlayerWithGinContext(c, auth_player)
	}
}

func playerAuth(token string) (gserv.AuthPlayer, error) {
	if token == "" {
		return gserv.AuthPlayer{}, nil
	}

	auth_player, err := jwt.ParseAuthPlayerToken(token)
	if err != nil {
		return gserv.AuthPlayer{}, err
	}

	return auth_player, nil
}

const (
	AUTH_PLAYER_KEY = "player"
)

func setAuthPlayerWithGinContext(c *gin.Context, player gserv.AuthPlayer) {
	c.Set(AUTH_PLAYER_KEY, player)
}

func GetAuthPlayerFromGinContext(c *gin.Context) gserv.AuthPlayer {
	player, _ := c.Get(AUTH_PLAYER_KEY)
	return player.(gserv.AuthPlayer)
}
