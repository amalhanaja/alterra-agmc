package controllers

import (
	"alterra-agmc-day-3/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetBooks(c echo.Context) error {
	books := []models.Book{
		{
			ID:        1234,
			Title:     "Buku Pertama",
			Isbn:      "1-234-5678-9101112-13",
			Writer:    "Alfian Akmal Hanantio",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
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
	book := models.Book{
		ID:        uint(id),
		Title:     "Buku Pertama",
		Isbn:      "1-234-5678-9101112-13",
		Writer:    "Alfian Akmal Hanantio",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   book,
	})
}

func CreateBook(c echo.Context) error {
	var jsonBody models.Book
	err := c.Bind(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
	}
	jsonBody.CreatedAt = time.Now().UTC()
	jsonBody.UpdatedAt = time.Now().UTC()
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "CREATED",
		"data":   jsonBody,
	})
}

func UpdateBook(c echo.Context) error {
	var jsonBody models.Book
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
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "OK",
		"data":   jsonBody,
	})
}

func DeleteBook(c echo.Context) error {
	strId := c.Param("id")
	_, err := strconv.Atoi(strId)
	if err != nil {
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
