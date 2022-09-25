package handlers

import (
	"alterra-agmc-day-7/internal/transportlayers/http/response"
	"alterra-agmc-day-7/pkg/jwt"
	"net/http"

	"github.com/labstack/echo/v4"

	goJWT "github.com/golang-jwt/jwt"
)

func getAuthorizedUserId(c echo.Context) (uint, error) {
	token, ok := c.Get("user").(*goJWT.Token)
	if !ok {
		return 0, c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Code:    "UNAUTHORIZED",
			Message: "failed get user",
		})
	}
	uid, err := jwt.ExtractID(token)
	if err != nil {
		return 0, c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Code:    "UNAUTHORIZED",
			Message: "invalid token",
		})
	}
	return uid, nil
}
