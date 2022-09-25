package handlers

import (
	"alterra-agmc-day-7/internal/models"
	"alterra-agmc-day-7/internal/services"
	"alterra-agmc-day-7/internal/transportlayers/http/request"
	"alterra-agmc-day-7/internal/transportlayers/http/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Login(c echo.Context) error
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

type userHandlerImpl struct {
	userService services.UserService
}

// Create implements UserHandler
func (h *userHandlerImpl) Create(c echo.Context) error {
	var requestBody request.CreateUserRequest
	err := c.Bind(&requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	if err := c.Validate(requestBody); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			})
		}
	}
	userToCreate := &models.User{
		Name:     requestBody.Name,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	createdUser, err := h.userService.Create(c.Request().Context(), userToCreate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	userResponse := response.UserResponse{
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}
	return c.JSON(http.StatusCreated, response.SuccessResponse[response.UserResponse]{
		Status: http.StatusCreated,
		Data:   userResponse,
	})
}

// Delete implements UserHandler
func (h *userHandlerImpl) Delete(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	if err := h.userService.DeleteByID(c.Request().Context(), uint(id), uid); err != nil {

		switch err := err.(type) {
		case services.ErrUnauthorized:
			return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Code:    "UNAUTHORIZED",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[any]{
		Status: http.StatusOK,
		Data:   nil,
	})
}

// GetAll implements UserHandler
func (h *userHandlerImpl) GetAll(c echo.Context) error {
	users, err := h.userService.FindAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	usersResponse := []response.UserResponse{}
	for _, u := range users {
		usersResponse = append(usersResponse, response.UserResponse{
			Name:      u.Name,
			Email:     u.Email,
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[[]response.UserResponse]{
		Status: http.StatusOK,
		Data:   usersResponse,
	})
}

// GetByID implements UserHandler
func (h *userHandlerImpl) GetByID(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	user, err := h.userService.FindByID(c.Request().Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	userResponse := response.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[response.UserResponse]{
		Status: http.StatusOK,
		Data:   userResponse,
	})
}

// Update implements UserHandler
func (h *userHandlerImpl) Update(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	var requestBody request.UpdateUserRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	if err := c.Validate(requestBody); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			})
		}
	}
	userToUpdate := &models.User{
		ID:       uint(id),
		Name:     requestBody.Name,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	updatedUser, err := h.userService.Update(c.Request().Context(), userToUpdate, uid)
	if err != nil {
		switch err := err.(type) {
		case services.ErrUnauthorized:
			return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Code:    "UNAUTHORIZED",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			})
		}
	}
	userResponse := response.UserResponse{
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[response.UserResponse]{
		Status: http.StatusOK,
		Data:   userResponse,
	})
}

// Login implements UserHandler
func (h *userHandlerImpl) Login(c echo.Context) error {
	var requestBody request.LoginUserRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	if err := c.Validate(requestBody); err != nil {
		switch err := err.(type) {
		case *echo.HTTPError:
			return c.JSON(err.Code, err.Message)
		default:
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			})
		}
	}
	token, err := h.userService.Login(c.Request().Context(), requestBody.Email, requestBody.Password)
	if err != nil {
		switch err := err.(type) {
		case services.ErrUnauthorized:
			return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Code:    "UNAUTHORIZED",
				Message: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_SERVER_ERROR",
				Message: err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[response.LoginResponse]{
		Status: http.StatusOK,
		Data: response.LoginResponse{
			Token: token,
		},
	})
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandlerImpl{
		userService: userService,
	}
}
