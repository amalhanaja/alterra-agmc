package controllers

import (
	"alterra-agmc-day-3/lib/jwt"
	"net/http"

	"github.com/labstack/echo/v4"

	goJWT "github.com/golang-jwt/jwt"
)

func getAuthorizedUserId(c echo.Context) (uint, error) {
	token, ok := c.Get("user").(*goJWT.Token)
	if !ok {
		return 0, c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "failed to get user",
		})
	}
	uid, err := jwt.ExtractID(token)
	if err != nil {
		return 0, c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "INTERNAL_SERVER_ERROR",
			"message": "failed to extract",
		})
	}
	return uid, nil
}
