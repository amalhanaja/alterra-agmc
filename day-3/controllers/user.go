package controllers

import (
	"alterra-agmc-day-3/lib/database"
	"alterra-agmc-day-3/lib/jwt"
	"alterra-agmc-day-3/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func LoginUser(c echo.Context) error {
	var jsonBody models.User
	if err := c.Bind(&jsonBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	user, err := database.GetUserByEmail(jsonBody.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": err.Error(),
		})
	}
	if user.Password != jsonBody.Password {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "username and password is not match",
		})
	}
	token, err := jwt.NewToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

func GetUsers(c echo.Context) error {
	users, err := database.GetUsers()
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "OK",
			"data":    []interface{}{},
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   users,
	})
}

func GetUserById(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	user, err := database.GetUserById(uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   user,
	})
}

func CreateUser(c echo.Context) error {
	var jsonBody models.User
	err := c.Bind(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	createdUser, err := database.CreateUser(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   createdUser,
	})
}

func UpdateUser(c echo.Context) error {
	var jsonBody models.User
	err := c.Bind(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	jsonBody.ID = uint(id)
	updatedUser, err := database.UpdateUser(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   updatedUser,
	})
}

func DeleteUser(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	if err := database.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "OK",
		"message": "deleted",
	})
}
