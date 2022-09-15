package middlewares

import (
	"alterra-agmc-day-3/config"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWT() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(config.GetJWTSecretKey()),
		SigningMethod: jwt.SigningMethodHS256.Name,
	})
}
