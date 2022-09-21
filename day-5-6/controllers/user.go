package controllers

import (
	"alterra-agmc-day-5-6/lib/database"
	"alterra-agmc-day-5-6/models"
	"alterra-agmc-day-5-6/pkg/jwt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func LoginUser(c echo.Context) error {
	var payload models.LoginUserPaylaod
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	if err := c.Validate(payload); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, "unknwon")
		}
	}
	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": err.Error(),
		})
	}
	if user.Password != payload.Password {
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
	var payload models.CreateUserPayload
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	if err := c.Validate(payload); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, "unknwon")
		}
	}
	user := &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}
	createdUser, err := database.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "CREATED",
		"data":   createdUser,
	})
}

func UpdateUser(c echo.Context) error {
	var payload models.UpdateUserPayload
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	if err := c.Validate(payload); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, "unknwon")
		}
	}
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	if uint(id) != uid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "access denied",
		})
	}
	user := &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}
	user.ID = uint(uid)
	updatedUser, err := database.UpdateUser(user)
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
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}

	if uint(id) != uid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "access denied",
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
