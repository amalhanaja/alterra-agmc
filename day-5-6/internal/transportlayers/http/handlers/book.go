package handlers

import (
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/services"
	"alterra-agmc-day-5-6/internal/transportlayers/http/request"
	"alterra-agmc-day-5-6/internal/transportlayers/http/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookHandler interface {
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Delete(c echo.Context) error
	Update(c echo.Context) error
	Create(c echo.Context) error
}

type bookHandlerImpl struct {
	service services.BookService
}

// Create implements BookHandler
func (h *bookHandlerImpl) Create(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	var requestBody request.CreateBookRequest
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
	book, err := h.service.Create(
		c.Request().Context(),
		&models.Book{
			Title:  requestBody.Title,
			Isbn:   requestBody.Isbn,
			Writer: requestBody.Writer,
			UserID: uid,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	bookResponse := response.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Writer:    book.Writer,
		Isbn:      book.Isbn,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse[response.BookResponse]{
		Status: http.StatusCreated,
		Data:   bookResponse,
	})
}

// Delete implements BookHandler
func (h *bookHandlerImpl) Delete(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	bookIdParam := c.Param("id")
	bookId, err := strconv.Atoi(bookIdParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	if err := h.service.DeleteByID(c.Request().Context(), uid, uint(bookId)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[any]{
		Status: http.StatusOK,
		Data:   nil,
	})
}

// GetAll implements BookHandler
func (h *bookHandlerImpl) GetAll(c echo.Context) error {
	books, err := h.service.FindAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	booksResponse := []response.BookResponse{}
	for _, b := range books {
		bookResponse := response.BookResponse{
			ID:        b.ID,
			Title:     b.Title,
			Writer:    b.Writer,
			Isbn:      b.Isbn,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
		booksResponse = append(booksResponse, bookResponse)
	}
	return c.JSON(http.StatusOK, response.SuccessResponse[[]response.BookResponse]{
		Status: http.StatusOK,
		Data:   booksResponse,
	})
}

// GetByID implements BookHandler
func (h *bookHandlerImpl) GetByID(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
	}
	book, err := h.service.FindByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
	}

	bookResponse := response.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Writer:    book.Writer,
		Isbn:      book.Isbn,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse[response.BookResponse]{
		Status: http.StatusOK,
		Data:   bookResponse,
	})
}

// Update implements BookHandler
func (h *bookHandlerImpl) Update(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	var requestBody request.UpdateBookRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
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
	book, err := h.service.Update(
		c.Request().Context(),
		&models.Book{
			ID:     uint(id),
			Title:  requestBody.Title,
			Isbn:   requestBody.Isbn,
			Writer: requestBody.Writer,
			UserID: uid,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
	bookResponse := response.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Writer:    book.Writer,
		Isbn:      book.Isbn,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse[response.BookResponse]{
		Status: http.StatusOK,
		Data:   bookResponse,
	})
}

func NewBookHandler(bookService services.BookService) BookHandler {
	return &bookHandlerImpl{
		service: bookService,
	}
}
