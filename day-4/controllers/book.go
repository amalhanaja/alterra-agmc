package controllers

import (
	"alterra-agmc-day-4/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var books []models.Book = make([]models.Book, 0)

func GetBooks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   books,
	})
}

func GetBookById(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	book := firstBook(func(b models.Book) bool {
		return b.ID == uint(id)
	})
	if book == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": "book not found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   book,
	})
}

func CreateBook(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	var payload models.CreateBookPayload
	utcNow := time.Now().UTC()
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
	book := models.Book{
		ID:        uint(len(books) + 1),
		Title:     payload.Title,
		Isbn:      payload.Isbn,
		Writer:    payload.Writer,
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
		UserID:    uid,
	}
	books = append(books, book)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "CREATED",
		"data":   book,
	})
}

func UpdateBook(c echo.Context) error {
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	var payload models.UpdateBookPayload
	if err := c.Bind(&payload); err != nil {
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
	if err := c.Validate(payload); err != nil {
		return err
	}
	book := firstBook(func(b models.Book) bool {
		return b.ID == uint(id)
	})
	if book == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": "book not found",
		})
	}
	if book.UserID != uid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "akses di tolak",
		})
	}
	if payload.Title != "" {
		book.Title = payload.Title
	}
	if payload.Isbn != "" {
		book.Isbn = payload.Isbn
	}
	if payload.Writer != "" {
		book.Writer = payload.Writer
	}
	books[id-1] = *book
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   book,
	})
}

func DeleteBook(c echo.Context) error {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	uid, err := getAuthorizedUserId(c)
	if err != nil {
		return err
	}
	book := firstBook(func(b models.Book) bool {
		return b.ID == uint(id)
	})
	if book == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": "book not found",
		})
	}
	if book.UserID != uid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  "UNAUTHORIZED",
			"message": "akses di tolak",
		})
	}
	books = append(books[:id-1], books[id:]...)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "OK",
		"message": "deleted",
	})
}

func firstBook(predicate func(models.Book) bool) *models.Book {
	for _, book := range books {
		if predicate(book) {
			return &book
		}
	}
	return nil
}
