package logic

import (
	"time"

	"demacia/common/baseauth"
	"demacia/service/auth/api/internal/types"

	"github.com/dgrijalva/jwt-go"
)

func buildTokens(authConfig baseauth.AuthConfig, fields map[string]interface{}) (types.Token, error) {
	var tokens types.Token

	accessToken, err := genToken(authConfig.AccessSecret, fields, authConfig.AccessExpire)
	if err != nil {
		return tokens, err
	}

	refreshToken, err := genToken(authConfig.RefreshSecret, fields, authConfig.RefreshExpire)
	if err != nil {
		return tokens, err
	}

	now := time.Now().Unix()
	tokens.AccessToken = accessToken
	tokens.AccessExpire = now + authConfig.AccessExpire
	tokens.RefreshAfter = now + authConfig.RefreshAfter
	tokens.RefreshToken = refreshToken
	tokens.RefreshExpire = now + authConfig.RefreshExpire
	return tokens, nil
}

func genToken(secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
