package controllers

import (
	"alterra-agmc-day-4/config"
	"alterra-agmc-day-4/lib/database"
	"alterra-agmc-day-4/lib/validator"
	"alterra-agmc-day-4/models"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setUp()
	returnCode := m.Run()
	tearDown()
	os.Exit(returnCode)
}

func setUp() {
	config.InitDB()
	database.CreateUser(&models.User{
		Email:    "user_1@email.com",
		Password: "Password",
		Name:     "user_1",
	})
}

func tearDown() {
	if err := config.DB.Unscoped().Delete(&models.User{}, "true").Error; err != nil {
		log.Fatal(err)
	}
	if err := config.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1").Error; err != nil {
		log.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name            string
		payload         map[string]interface{}
		expectedCode    int
		expectedStatus  string
		expectedMessage *struct{ value string }
	}{
		{
			name:            "Test Login when payload is invalid should return bad request with message",
			payload:         map[string]interface{}{},
			expectedCode:    http.StatusBadRequest,
			expectedStatus:  "BAD_REQUEST",
			expectedMessage: &struct{ value string }{"Key: 'LoginUserPaylaod.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginUserPaylaod.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		},
		{
			name: "Teset login when payload is valid and user not found should return unauthorized",
			payload: map[string]interface{}{
				"email":    "user_not_found@email.com",
				"password": "Password",
			},
			expectedCode:    http.StatusUnauthorized,
			expectedStatus:  "UNAUTHORIZED",
			expectedMessage: &struct{ value string }{"record not found"},
		},
		{
			name: "Teset login when payload is valid and password is not match should return unauthorized",
			payload: map[string]interface{}{
				"email":    "user_1@email.com",
				"password": "Incorrect",
			},
			expectedCode:    http.StatusUnauthorized,
			expectedStatus:  "UNAUTHORIZED",
			expectedMessage: &struct{ value string }{"username and password is not match"},
		},
		{
			name: "Test login when valid, user is found, and  password is correct should return ok with token payload",
			payload: map[string]interface{}{
				"email":    "user_1@email.com",
				"password": "Password",
			},
			expectedCode:    http.StatusOK,
			expectedStatus:  "OK",
			expectedMessage: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			jsonPayload, _ := json.Marshal(&testCase.payload)
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/login")

			// Act
			LoginUser(c)

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
		name            string
		expectedCode    int
		expectedStatus  string
		expectedMessage *struct{ value string }
	}{
		{
			name:           "Test get users should return ok with data",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
		},
	}

	for _, testCase := range testCases {
		// Arrange
		e := echo.New()
		e.Validator = validator.NewCustomValidator()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users")

		// Act
		GetUsers(c)

		// Assert
		var payload map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedCode, rec.Code)
		if testCase.expectedMessage != nil {
			assert.Equal(t, testCase.expectedMessage.value, payload["message"])
		}
		if payload["data"] != nil {
			data := payload["data"]
			assert.NotEmpty(t, data)
		}
	}
}

func TestGetuserById(t *testing.T) {
	testCases := []struct {
		name            string
		userId          string
		expectedCode    int
		expectedStatus  string
		expectedMessage *struct{ value string }
		expectedData    *models.User
	}{
		{
			name:            "Test get user by id when id is NaN should return bad request with message",
			userId:          "NaN",
			expectedCode:    http.StatusBadRequest,
			expectedStatus:  "BAD_REQUEST",
			expectedMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name:            "Test get user by id when not found should return bad request with message",
			userId:          "1002",
			expectedCode:    http.StatusBadRequest,
			expectedStatus:  "BAD_REQUEST",
			expectedMessage: &struct{ value string }{"record not found"},
		},
		{
			name:           "Test get user by id when user is found should return ok with user data",
			userId:         "1",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectedData: &models.User{
				Name:  "user_1",
				Email: "user_1@email.com",
			},
		},
	}

	for _, testCase := range testCases {
		// Arrange
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
		GetUserById(c)

		// Assert
		var payload map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, testCase.expectedStatus, payload["status"])
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
		name            string
		userPayload     map[string]interface{}
		expectedCode    int
		expectedStatus  string
		expectedMessage *struct{ value string }
		expectedData    *models.User
	}{
		{
			name:            "Test create user when user is invalid should return bad request with message",
			userPayload:     map[string]interface{}{},
			expectedCode:    http.StatusBadRequest,
			expectedStatus:  "BAD_REQUEST",
			expectedMessage: &struct{ value string }{"Key: 'CreateUserPayload.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'CreateUserPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'CreateUserPayload.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		},
		{
			name: "Test create user when user is valid should return bad created with data",
			userPayload: map[string]interface{}{
				"email":    "user_2@email.com",
				"name":     "user_2",
				"password": "secret_password",
			},
			expectedCode:   http.StatusCreated,
			expectedStatus: "CREATED",
			expectedData: &models.User{
				Email: "user_2@email.com",
				Name:  "user_2",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			jsonPayload, _ := json.Marshal(&testCase.userPayload)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users")

			// Act
			CreateUser(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
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
		name            string
		userId          string
		userPayload     map[string]interface{}
		token           *jwt.Token
		expectedCode    int
		expectedStatus  string
		expectedMessage *struct{ value string }
		expectedData    *models.User
	}{
		{
			name:            "Test update user when user is unauthorized should return unauthorized with message",
			userPayload:     map[string]interface{}{},
			expectedCode:    http.StatusUnauthorized,
			expectedStatus:  "UNAUTHORIZED",
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
			expectedStatus:  "BAD_REQUEST",
			expectedMessage: &struct{ value string }{`strconv.Atoi: parsing "NaN": invalid syntax`},
		},
		{
			name: "Test update user when user is authorized and doesn't has access to edit user",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "100"},
			},
			userId:          "1",
			expectedCode:    http.StatusUnauthorized,
			expectedStatus:  "UNAUTHORIZED",
			expectedMessage: &struct{ value string }{"access denied"},
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
			userId:         "1",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectedData: &models.User{
				Name:  "New User",
				Email: "new_email@email.com",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
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
			UpdateUser(c)

			// Assert
			var payload map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&payload)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, rec.Code)
			assert.Equal(t, testCase.expectedStatus, payload["status"])
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
		name           string
		userId         string
		token          *jwt.Token
		expectedCode   int
		expectedStatus string
		expectMessage  string
	}{
		{
			name:           "Test delete user when user is unauthorized should return unauthorized with message",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  "failed get user",
		},
		{
			name:   "Test delete user when id is NaN should return bad request with message",
			userId: "NaN",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "123"},
			},
			expectedCode:   http.StatusBadRequest,
			expectedStatus: "BAD_REQUEST",
			expectMessage:  `strconv.Atoi: parsing "NaN": invalid syntax`,
		},
		{
			name: "Test delete user when user is authorized and doesn't has access to delete user",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			userId:         "123",
			expectedCode:   http.StatusUnauthorized,
			expectedStatus: "UNAUTHORIZED",
			expectMessage:  "access denied",
		},
		{
			name: "Test update book when user is authorized and has access to edit book",
			token: &jwt.Token{
				Valid:  true,
				Claims: jwt.MapClaims{"sub": "1"},
			},
			userId:         "1",
			expectedCode:   http.StatusOK,
			expectedStatus: "OK",
			expectMessage:  "deleted",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
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
			DeleteUser(c)

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
