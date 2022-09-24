package handlers

import (
	"alterra-agmc-day-5-6/internal/datasources"
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/services/mocks"
	"alterra-agmc-day-5-6/pkg/validator"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllBooks(t *testing.T) {
	testCases := []struct {
		name                   string
		booksReturnFromService []*models.Book
		errReturnFromService   error
		expectedCode           int
		expectedStatus         string
		expectBooksCount       int
	}{
		{
			name:                   "Test GetAll when books is empty should return ok with empty data",
			booksReturnFromService: make([]*models.Book, 0),
			expectedCode:           http.StatusOK,
			expectedStatus:         "OK",
			expectBooksCount:       0,
		},
		{
			name: "Test GetAll when books not empty should return ok books data",
			booksReturnFromService: []*models.Book{
				{
					ID:        123,
					Title:     "Test Book",
					Isbn:      "ISBN",
					Writer:    "Alfian Akmal Hanantio",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					UserID:    1001,
				},
			},
			expectedCode:     http.StatusOK,
			expectedStatus:   "OK",
			expectBooksCount: 1,
		},
		{
			name:                   "Test GetAll when service error should return internal service error",
			booksReturnFromService: nil,
			errReturnFromService:   errors.New("something bad"),
			expectedCode:           500,
			expectBooksCount:       0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			bookService := mocks.NewBookService(t)
			bookHandler := NewBookHandler(bookService)
			bookService.On("FindAll", mock.Anything).Return(testCase.booksReturnFromService, testCase.errReturnFromService)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Act
			result := bookHandler.GetAll(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, result)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.errReturnFromService == nil {
				assert.Equal(t, testCase.expectBooksCount, len(payload["data"].([]interface{})))
			}
		})
	}
}

func TestGetBookById(t *testing.T) {
	testCases := []struct {
		name          string
		bookId        string
		bookReturn    *models.Book
		errReturn     error
		expectedCode  int
		expectMessage *struct{ value string }
		expectBook    *models.Book
	}{
		{
			name:          "Test get book by id when id is NaN should return bad request with message",
			bookId:        "NaN",
			expectedCode:  http.StatusBadRequest,
			expectMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:          "Test get book by id when book is not found should return bad requst with message",
			bookId:        "123",
			errReturn:     datasources.ErrRecordNotFound{},
			expectedCode:  http.StatusInternalServerError,
			expectMessage: &struct{ value string }{"record not found"},
		},
		{
			name:   "Test get book by id when book is found should return ok with book data",
			bookId: "100",
			bookReturn: &models.Book{
				ID:        100,
				Title:     "Test Book",
				Isbn:      "ISBN",
				Writer:    "Alfian Akmal Hanantio",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				UserID:    1001,
			},
			expectedCode:  http.StatusOK,
			expectMessage: nil,
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
			bookService := mocks.NewBookService(t)
			bookHandler := NewBookHandler(bookService)
			if testCase.bookReturn != nil || testCase.errReturn != nil {
				bookService.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(testCase.bookReturn, testCase.errReturn)
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.bookId)

			// Act
			result := bookHandler.GetByID(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, result)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
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
		name          string
		bookPayload   map[string]interface{}
		bookReturn    *models.Book
		errReturn     error
		token         *jwt.Token
		expectedCode  int
		expectMessage *struct{ value string }
		expectBook    *models.Book
	}{
		{
			name:          "Test create book when user is unauthorized should return unauthorized with message",
			bookPayload:   map[string]interface{}{},
			expectedCode:  http.StatusUnauthorized,
			expectMessage: &struct{ value string }{"failed get user"},
		},
		{
			name:        "Test create book when user is authorized and book is invalid should return bad request with message",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:  http.StatusBadRequest,
			expectMessage: &struct{ value string }{"Key: 'CreateBookRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'CreateBookRequest.Isbn' Error:Field validation for 'Isbn' failed on the 'required' tag\nKey: 'CreateBookRequest.Writer' Error:Field validation for 'Writer' failed on the 'required' tag"},
		},
		{
			name: "Test create book when user is authorized and service error should reteurn internal server error",
			bookPayload: map[string]interface{}{
				"title":  "Test Book",
				"writer": "Alfian Akmal Hanantio",
				"isbn":   "ISBN",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			errReturn:     errors.New("unknown error"),
			expectedCode:  http.StatusInternalServerError,
			expectMessage: &struct{ value string }{"unknown error"},
		},
		{
			name: "Test create book when user is authorized and book is valid should return bad created with data",
			bookPayload: map[string]interface{}{
				"title":  "Test Book",
				"writer": "Alfian Akmal Hanantio",
				"isbn":   "ISBN",
			},
			bookReturn: &models.Book{
				ID:     100,
				Title:  "Test Book",
				Isbn:   "ISBN",
				Writer: "Alfian Akmal Hanantio",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode: http.StatusCreated,
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
			bookService := mocks.NewBookService(t)
			bookHandler := NewBookHandler(bookService)
			if testCase.bookReturn != nil || testCase.errReturn != nil {
				bookService.On("Create", mock.Anything, mock.Anything).Return(testCase.bookReturn, testCase.errReturn)
			}
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
			bookHandler.Create(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
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
		name          string
		bookId        string
		bookPayload   map[string]interface{}
		bookReturn    *models.Book
		errReturn     error
		token         *jwt.Token
		expectedCode  int
		expectMessage *struct{ value string }
		expectBook    *models.Book
	}{
		{
			name:          "Test update book when user is unauthorized should return unauthorized with message",
			bookPayload:   map[string]interface{}{},
			expectedCode:  http.StatusUnauthorized,
			expectMessage: &struct{ value string }{"failed get user"},
		},
		{
			name:   "Test update book when id is NaN should return bad request with message",
			bookId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:  http.StatusBadRequest,
			expectMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:        "Test update book when user is authorized and book not found",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			errReturn:     datasources.ErrRecordNotFound{},
			bookId:        "100",
			expectedCode:  http.StatusInternalServerError,
			expectMessage: &struct{ value string }{"record not found"},
		},
		{
			name:        "Test update book when user is authorized and return error from service should return internal server error",
			bookPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			errReturn:     errors.New("error"),
			bookId:        "123",
			expectedCode:  http.StatusInternalServerError,
			expectMessage: &struct{ value string }{"error"},
		},
		{
			name: "Test update book when user is authorized and success edit book should return ok",
			bookPayload: map[string]interface{}{
				"title":  "New Title",
				"isbn":   "New ISBN",
				"writer": "New Writer",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			bookReturn: &models.Book{
				ID:     1,
				Title:  "New Title",
				Isbn:   "New ISBN",
				Writer: "New Writer",
			},
			bookId:       "1",
			expectedCode: http.StatusOK,
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
			bookService := mocks.NewBookService(t)
			bookHandler := NewBookHandler(bookService)
			if testCase.bookReturn != nil || testCase.errReturn != nil {
				bookService.On("Update", mock.Anything, mock.Anything).Return(testCase.bookReturn, testCase.errReturn)
			}
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
			bookHandler.Update(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
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
		name          string
		bookId        string
		errReturn     error
		callService   bool
		token         *jwt.Token
		expectedCode  int
		expectMessage *struct{ value string }
	}{
		{
			name:          "Test delete book when user is unauthorized should return unauthorized with message",
			expectedCode:  http.StatusUnauthorized,
			expectMessage: &struct{ value string }{"failed get user"},
		},
		{
			name:   "Test delete book when id is NaN should return bad request with message",
			bookId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:  http.StatusBadRequest,
			expectMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name: "Test delete book when user is authorized and service error should return internal server error",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			callService:   true,
			errReturn:     errors.New("error"),
			bookId:        "100",
			expectedCode:  http.StatusInternalServerError,
			expectMessage: &struct{ value string }{"error"},
		},
		{
			name: "Test delete book when user is authorized and has access to delete book",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			callService:  true,
			bookId:       "1",
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			bookService := mocks.NewBookService(t)
			bookHandler := NewBookHandler(bookService)
			if testCase.callService {
				bookService.On("DeleteByID", mock.Anything, mock.Anything, mock.Anything).Return(testCase.errReturn)
			}
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
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
			bookHandler.Delete(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectMessage != nil {
				assert.Equal(t, testCase.expectMessage.value, payload["message"])
			}
		})
	}
}
