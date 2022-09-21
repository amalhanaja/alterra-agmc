package jwt

import (
	"alterra-agmc-day-5-6/config"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(id uint) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(id)),
		ExpiresAt: time.Now().Add(time.Millisecond * time.Duration(config.GetJWTExpirationTime())).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecretKey()))
}

func ExtractID(token *jwt.Token) (uint, error) {
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, fmt.Errorf("failed parse id")
	}
	return uint(id), nil
}
