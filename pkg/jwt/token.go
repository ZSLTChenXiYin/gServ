package jwt

import (
	"errors"
	"gServ/core/config"
	"gServ/pkg/gserv"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AUTH_PLAYER_KEY_ID         = "id"
	AUTH_PLAYER_KEY_EXPIRED_AT = "expired_at"
)

func CreateAuthPlayerToken(auth_player gserv.AuthPlayer) (string, error) {
	expiration_time := time.Now().Add(time.Hour * 24 * 7)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		AUTH_PLAYER_KEY_ID:         auth_player.ID,
		AUTH_PLAYER_KEY_EXPIRED_AT: expiration_time.Unix(),
	}).SignedString(config.GetConfig().Server.Jwt)
	return token, err
}

func ParseAuthPlayerToken(signed_token string) (gserv.AuthPlayer, error) {
	auth_player := gserv.AuthPlayer{}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(signed_token, claims, func(token *jwt.Token) (any, error) {
		return config.GetConfig().Server.Jwt, nil
	})
	if err != nil {
		return auth_player, err
	}
	if !token.Valid {
		return auth_player, errors.New("invalid token")
	}

	auth_player = gserv.AuthPlayer{
		ID:        uint((*claims)[AUTH_PLAYER_KEY_ID].(float64)),
		ExpiredAt: time.Unix(int64((*claims)[AUTH_PLAYER_KEY_EXPIRED_AT].(float64)), 0),
	}

	return auth_player, nil
}

func RefreshAuthPlayerToken(signed_token string) (string, error) {
	auth_player, err := ParseAuthPlayerToken(signed_token)
	if err != nil {
		return "", err
	}

	return CreateAuthPlayerToken(auth_player)
}

// GenerateToken 生成JWT令牌（兼容现有代码）
func GenerateToken(playerID uint) (string, error) {
	authPlayer := gserv.AuthPlayer{
		ID: playerID,
	}
	return CreateAuthPlayerToken(authPlayer)
}
