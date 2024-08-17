package jwt

import (
	"github.com/golang-jwt/jwt"
	"vk_tarantool_project/internal/domain"
)

func GetMapClaims(tokenString, secret string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
