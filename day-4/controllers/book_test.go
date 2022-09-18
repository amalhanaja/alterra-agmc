package controllers

import (
	"alterra-agmc-day-4/lib/validator"
	"alterra-agmc-day-4/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	testCases := []struct {
		name             string
		before           func()
		expectedCode     int
		expectedStatus   string
		expectBooksCount int
	}{
		{
			name: "Test Get Books when books is empty should return empty data",
			before: func() {
				books = make([]models.Book, 0)
			},
			expectedCode:     http.StatusOK,
			expectedStatus:   "OK",
			expectBooksCount: 0,
		},
		{
			name: "Test Get Books when books not empty should return books data",
			before: func() {
				books = []models.Book{
					{
						ID:        123,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    1001,
					},
				}
			},
			expectedCode:     http.StatusOK,
			expectedStatus:   "OK",
			expectBooksCount: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			testCase.before()
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Act
			result := GetBooks(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, result)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectBooksCount, len(payload["data"].([]interface{})))
		})
	}
}

func TestGetBookById(t *testing.T) {
	testCases := []struct {
		name           string
		bookId         string
		before         func()
		expectedCode   int
		expectedStatus string
		expectMessage  *struct{ value string }
		expectBook     *models.Book
	}{
		{
			name:   "Test get book by id when id is NaN should return bad request with message",
			bookId: "NaN",
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:   "Test get book by id when book is not found should return bad requst with message",
			bookId: "123",
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  &struct{ value string }{"book not found"},
		},
		{
			name:   "Test get book by id when book is found should return ok with book data",
			bookId: "100",
			before: func() {
				books = []models.Book{
					{
						ID:        100,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    1001,
					},
				}
			},
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectMessage:  nil,
			expectBook: &models.Book{
				ID:     100,
				Title:  "Test Book",
				Isbn:   "ISBN",
				Writer: "Alfian Akmal Hanantio",
				UserID: 1001,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			testCase.before()
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.bookId)

			// Act
			result := GetBookById(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, result)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
			if testCase.expectMessage != nil {
				assert.Equal(t, testCase.expectMessage.value, payload["message"])
			}
			if testCase.expectBook != nil {
				data := payload["data"].(map[string]interface{})
				assert.Equal(t, testCase.expectBook.ID, uint(data["id"].(float64)))
				assert.Equal(t, testCase.expectBook.Isbn, data["isbn"])
				assert.Equal(t, testCase.expectBook.Title, data["title"])
				assert.Equal(t, testCase.expectBook.Writer, data["writer"])
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	testCases := []struct {
		name           string
		bookPayload    map[string]interface{}
		token          *jwt.Token
		expectedCode   int
		expectedStatus string
		expectMessage  *struct{ value string }
		expectBook     *models.Book
	}{
		{
			name:           "Test create book when user is unauthorized should return unauthorized with message",
			bookPayload:    map[string]interface{}{},
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  &struct{ value string }{"failed get user"},
		},
		{
			name:        "Test create book when user is authorized and book is invalid should return bad request with message",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  &struct{ value string }{"Key: 'CreateBookPayload.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'CreateBookPayload.Isbn' Error:Field validation for 'Isbn' failed on the 'required' tag\nKey: 'CreateBookPayload.Writer' Error:Field validation for 'Writer' failed on the 'required' tag"},
		},
		{
			name: "Test create book when user is authorized and book is valid should return bad created with data",
			bookPayload: map[string]interface{}{
				"title":  "Test Book",
				"writer": "Alfian Akmal Hanantio",
				"isbn":   "ISBN",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:   http.StatusCreated,
			expectedStatus: "CREATED",
			expectBook: &models.Book{
				ID:     100,
				Title:  "Test Book",
				Isbn:   "ISBN",
				Writer: "Alfian Akmal Hanantio",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			jsonPayload, _ := json.Marshal(&testCase.bookPayload)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books")
			if testCase.token != nil {
				c.Set("user", testCase.token)
			}

			// Act
			CreateBook(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
			if testCase.expectMessage != nil {
				assert.Equal(t, testCase.expectMessage.value, payload["message"])
			}
			if testCase.expectBook != nil {
				data := payload["data"].(map[string]interface{})
				assert.Equal(t, testCase.expectBook.Isbn, data["isbn"])
				assert.Equal(t, testCase.expectBook.Title, data["title"])
				assert.Equal(t, testCase.expectBook.Writer, data["writer"])
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	testCases := []struct {
		name           string
		bookId         string
		bookPayload    map[string]interface{}
		before         func()
		token          *jwt.Token
		expectedCode   int
		expectedStatus string
		expectMessage  *struct{ value string }
		expectBook     *models.Book
	}{
		{
			name:        "Test update book when user is unauthorized should return unauthorized with message",
			bookPayload: map[string]interface{}{},
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  &struct{ value string }{"failed get user"},
		},
		{
			name:   "Test update book when id is NaN should return bad request with message",
			bookId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:        "Test update book when user is authorized and book not found",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{}
			},
			bookId:         "100",
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  &struct{ value string }{"book not found"},
		},
		{
			name:        "Test update book when user is authorized and doesn't has access to edit book",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{
					{
						ID:        123,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    1001,
					},
				}
			},
			bookId:         "123",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  &struct{ value string }{"access denied"},
		},
		{
			name: "Test update book when user is authorized and has access to edit book",
			bookPayload: map[string]interface{}{
				"title":  "New Title",
				"isbn":   "New ISBN",
				"writer": "New Writer",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{
					{
						ID:        1,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    123,
					},
				}
			},
			bookId:         "1",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectBook: &models.Book{
				ID:     1,
				Title:  "New Title",
				Isbn:   "New ISBN",
				Writer: "New Writer",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			testCase.before()
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			jsonPayload, _ := json.Marshal(&testCase.bookPayload)
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.bookId)
			if testCase.token != nil {
				c.Set("user", testCase.token)
			}

			// Act
			UpdateBook(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
			if testCase.expectMessage != nil {
				assert.Equal(t, testCase.expectMessage.value, payload["message"])
			}
			if testCase.expectBook != nil {
				data := payload["data"].(map[string]interface{})
				assert.Equal(t, testCase.expectBook.Isbn, data["isbn"])
				assert.Equal(t, testCase.expectBook.Title, data["title"])
				assert.Equal(t, testCase.expectBook.Writer, data["writer"])
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	testCases := []struct {
		name           string
		bookId         string
		before         func()
		token          *jwt.Token
		expectedCode   int
		expectedStatus string
		expectMessage  string
	}{
		{
			name: "Test delete book when user is unauthorized should return unauthorized with message",
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  "failed get user",
		},
		{
			name:   "Test delete book when id is NaN should return bad request with message",
			bookId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{}
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  `strconv.Atoi: parsing "NaN": invalid syntax`,
		},
		{
			name: "Test update book when user is authorized and book not found",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{}
			},
			bookId:         "100",
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  "book not found",
		},
		{
			name: "Test update book when user is authorized and doesn't has access to edit book",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{
					{
						ID:        123,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    1001,
					},
				}
			},
			bookId:         "123",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  "access denied",
		},
		{
			name: "Test update book when user is authorized and has access to edit book",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			before: func() {
				books = []models.Book{
					{
						ID:        1,
						Title:     "Test Book",
						Isbn:      "ISBN",
						Writer:    "Alfian Akmal Hanantio",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						UserID:    123,
					},
				}
			},
			bookId:         "1",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectMessage:  "deleted",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			testCase.before()
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.bookId)
			if testCase.token != nil {
				c.Set("user", testCase.token)
			}

			// Act
			DeleteBook(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
			assert.Equal(t, testCase.expectMessage, payload["message"])
		})
	}
}
