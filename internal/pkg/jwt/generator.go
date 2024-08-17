package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
	"vk_tarantool_project/internal/domain"
)

func NewToken(user domain.UserInfo, secret string, dur time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(dur).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
