package handlers

import (
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/services"
	"alterra-agmc-day-5-6/internal/services/mocks"
	"alterra-agmc-day-5-6/pkg/validator"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name             string
		payload          map[string]interface{}
		expectedCode     int
		errReturnService error
		returnService    string
		expectedMessage  *struct{ value string }
	}{
		{
			name:            "Test Login when request is invalid should return bad request with message",
			payload:         map[string]interface{}{},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: &struct{ value string }{"Key: 'LoginUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		},
		{
			name: "Test login when service return unauthorized should return unauthorized",
			payload: map[string]interface{}{
				"email":    "user_not_found@email.com",
				"password": "Password",
			},
			expectedCode:     http.StatusUnauthorized,
			errReturnService: services.ErrUnauthorized{},
			expectedMessage:  &struct{ value string }{"unauthorized"},
		},
		{
			name: "Test login when service err should return internal server error",
			payload: map[string]interface{}{
				"email":    "user_1@email.com",
				"password": "Incorrect",
			},
			errReturnService: errors.New("error"),
			expectedCode:     http.StatusInternalServerError,
			expectedMessage:  &struct{ value string }{"error"},
		},
		{
			name: "Test login when service return token should return ok",
			payload: map[string]interface{}{
				"email":    "user_1@email.com",
				"password": "Password",
			},
			returnService:   "TOKEN",
			expectedCode:    http.StatusOK,
			expectedMessage: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewUserService(t)
			handler := NewUserHandler(mockService)
			jsonPayload, _ := json.Marshal(&testCase.payload)
			if testCase.errReturnService != nil || testCase.returnService != "" {
				mockService.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(testCase.returnService, testCase.errReturnService)
			}
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/login")

			// Act
			handler.Login(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectedMessage != nil {
				assert.Equal(t, testCase.expectedMessage.value, payload["message"])
			}
			if payload["data"] != nil {
				data := payload["data"].(map[string]interface{})
				assert.NotNil(t, data["token"])
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	testCases := []struct {
		name             string
		expectedCode     int
		serviceReturn    []*models.User
		errServiceReturn error
		expectedMessage  *struct{ value string }
	}{
		{
			name:          "Test get users should return ok with data",
			expectedCode:  http.StatusOK,
			serviceReturn: make([]*models.User, 0),
		},
		{
			name:             "Test get all users when service return error should return internal server error",
			expectedCode:     http.StatusInternalServerError,
			errServiceReturn: errors.New("error"),
			expectedMessage:  &struct{ value string }{"error"},
		},
	}

	for _, testCase := range testCases {
		// Arrange
		mockService := mocks.NewUserService(t)
		handler := NewUserHandler(mockService)
		mockService.On("FindAll", mock.Anything).Return(testCase.serviceReturn, testCase.errServiceReturn)
		e := echo.New()
		e.Validator = validator.NewCustomValidator()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")

		// Act
		handler.GetAll(c)

		// Assert
		var payload map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedCode, rec.Code)
		if testCase.expectedMessage != nil {
			assert.Equal(t, testCase.expectedMessage.value, payload["message"])
		}
	}
}

func TestGetuserById(t *testing.T) {
	testCases := []struct {
		name             string
		userId           string
		expectedCode     int
		serviceReturn    *models.User
		errServiceReturn error
		expectedMessage  *struct{ value string }
		expectedData     *models.User
	}{
		{
			name:            "Test get user by id when id is NaN should return bad request with message",
			userId:          "NaN",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:             "Test get user by id service err should return internal server error",
			userId:           "1002",
			errServiceReturn: errors.New("error"),
			serviceReturn:    &models.User{},
			expectedCode:     http.StatusInternalServerError,
			expectedMessage:  &struct{ value string }{"error"},
		},
		{
			name:   "Test get user by id when user is found should return ok with user data",
			userId: "1",
			serviceReturn: &models.User{
				Name:  "user_1",
				Email: "user_1@email.com",
			},
			expectedCode: http.StatusOK,
			expectedData: &models.User{
				Name:  "user_1",
				Email: "user_1@email.com",
			},
		},
	}

	for _, testCase := range testCases {
		// Arrange
		mockService := mocks.NewUserService(t)
		handler := NewUserHandler(mockService)
		if testCase.errServiceReturn != nil || testCase.serviceReturn != nil {
			mockService.On("FindByID", mock.Anything, mock.Anything).Return(testCase.serviceReturn, testCase.errServiceReturn)
		}
		e := echo.New()
		e.Validator = validator.NewCustomValidator()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(testCase.userId)

		// Act
		handler.GetByID(c)

		// Assert
		var payload map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedCode, rec.Code)
		if testCase.expectedMessage != nil {
			assert.Equal(t, testCase.expectedMessage.value, payload["message"])
		}
		if testCase.expectedData != nil {
			data := payload["data"].(map[string]interface{})
			assert.Equal(t, testCase.expectedData.Email, data["email"])
			assert.Equal(t, testCase.expectedData.Name, data["name"])
		}
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name             string
		userPayload      map[string]interface{}
		serviceReturn    *models.User
		errServiceReturn error
		expectedCode     int
		expectedMessage  *struct{ value string }
		expectedData     *models.User
	}{
		{
			name:            "Test create user when user is invalid should return bad request with message",
			userPayload:     map[string]interface{}{},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: &struct{ value string }{"Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		}, {
			name: "Test create user when service error should return internal serverError",
			userPayload: map[string]interface{}{
				"email":    "user_2@email.com",
				"name":     "user_2",
				"password": "secret_password",
			},
			serviceReturn:    &models.User{},
			expectedCode:     http.StatusInternalServerError,
			errServiceReturn: errors.New("error"),
			expectedMessage:  &struct{ value string }{"error"},
		},
		{
			name: "Test create user when user is valid should return bad created with data",
			userPayload: map[string]interface{}{
				"email":    "user_2@email.com",
				"name":     "user_2",
				"password": "secret_password",
			},
			serviceReturn: &models.User{
				Email: "user_2@email.com",
				Name:  "user_2",
			},
			expectedCode: http.StatusCreated,
			expectedData: &models.User{
				Email: "user_2@email.com",
				Name:  "user_2",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewUserService(t)
			handler := NewUserHandler(mockService)
			if testCase.errServiceReturn != nil || testCase.serviceReturn != nil {
				mockService.On("Create", mock.Anything, mock.Anything).Return(testCase.serviceReturn, testCase.errServiceReturn)
			}
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			jsonPayload, _ := json.Marshal(&testCase.userPayload)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users")

			// Act
			handler.Create(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectedMessage != nil {
				assert.Equal(t, testCase.expectedMessage.value, payload["message"])
			}
			if testCase.expectedData != nil {
				data := payload["data"].(map[string]interface{})
				assert.Equal(t, testCase.expectedData.Email, data["email"])
				assert.Equal(t, testCase.expectedData.Name, data["name"])
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name             string
		userId           string
		userPayload      map[string]interface{}
		token            *jwt.Token
		serviceReturn    *models.User
		errServiceReturn error
		expectedCode     int
		expectedMessage  *struct{ value string }
		expectedData     *models.User
	}{
		{
			name:            "Test update user when user is unauthorized should return unauthorized with message",
			userPayload:     map[string]interface{}{},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: &struct{ value string }{"failed get user"},
		},
		{
			name:        "Test update user when id is NaN should return bad request with message",
			userId:      "NaN",
			userPayload: map[string]interface{}{},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name: "Test update user when user is authorized service return ErrUnauthorized should return unauthorized",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "100"},
			},
			userId:           "1",
			serviceReturn:    &models.User{},
			errServiceReturn: services.ErrUnauthorized{},
			expectedCode:     http.StatusUnauthorized,
			expectedMessage:  &struct{ value string }{"unauthorized"},
		},
		{
			name: "Test update user when user is authorized service error should return internal server error",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "100"},
			},
			userId:           "1",
			serviceReturn:    &models.User{},
			errServiceReturn: errors.New("error"),
			expectedCode:     http.StatusInternalServerError,
			expectedMessage:  &struct{ value string }{"error"},
		},
		{
			name: "Test update book when user is authorized and has access to edit book",
			userPayload: map[string]interface{}{
				"email": "new_email@email.com",
				"name":  "New User",
			},
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			userId: "1",
			serviceReturn: &models.User{
				Name:  "New User",
				Email: "new_email@email.com",
			},
			expectedCode: http.StatusOK,
			expectedData: &models.User{
				Name:  "New User",
				Email: "new_email@email.com",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewUserService(t)
			handler := NewUserHandler(mockService)
			if testCase.errServiceReturn != nil || testCase.serviceReturn != nil {
				mockService.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(testCase.serviceReturn, testCase.errServiceReturn)
			}
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			jsonPayload, _ := json.Marshal(&testCase.userPayload)
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/books/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.userId)
			if testCase.token != nil {
				c.Set("user", testCase.token)
			}

			// Act
			handler.Update(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectedMessage != nil {
				assert.Equal(t, testCase.expectedMessage.value, payload["message"])
			}
			if testCase.expectedData != nil {
				data := payload["data"].(map[string]interface{})
				assert.Equal(t, testCase.expectedData.Email, data["email"])
				assert.Equal(t, testCase.expectedData.Name, data["name"])
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name            string
		userId          string
		token           *jwt.Token
		expectedCode    int
		errReturn       error
		callService     bool
		expectedMessage *struct{ value string }
	}{
		{
			name:            "Test delete user when user is unauthorized should return unauthorized with message",
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: &struct{ value string }{"failed get user"},
		},
		{
			name:   "Test delete user when id is NaN should return bad request with message",
			userId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name: "Test delete user when user is authorized and doesn't has access to delete user",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			userId:          "123",
			errReturn:       services.ErrUnauthorized{},
			callService:     true,
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: &struct{ value string }{"unauthorized"},
		},
		{
			name: "Test delete user when service return error should return internal server error",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			userId:          "123",
			errReturn:       errors.New("error"),
			callService:     true,
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: &struct{ value string }{"error"},
		},
		{
			name: "Test update book when user is authorized and has access to edit book",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			callService:  true,
			errReturn:    nil,
			userId:       "1",
			expectedCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockService := mocks.NewUserService(t)
			handler := NewUserHandler(mockService)
			if testCase.callService {
				mockService.On("DeleteByID", mock.Anything, mock.Anything, mock.Anything).Return(testCase.errReturn)
			}
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id")
			c.SetParamNames("id")
			c.SetParamValues(testCase.userId)
			if testCase.token != nil {
				c.Set("user", testCase.token)
			}

			// Act
			handler.Delete(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			if testCase.expectedMessage != nil {
				assert.Equal(t, testCase.expectedMessage.value, payload["message"])
			}
		})
	}
}
