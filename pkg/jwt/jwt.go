package jwt

import (
	"time"

	"github.com/warthog618/modem/pkg/config"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(number string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    number,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.TomlConf.JWT.EXPIRES_AT)).Unix(),
	})
	return claims.SignedString([]byte(config.TomlConf.JWT.SECRETKEY))
}
