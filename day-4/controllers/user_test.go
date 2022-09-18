package controllers

import (
	"alterra-agmc-day-4/config"
	"alterra-agmc-day-4/lib/database"
	"alterra-agmc-day-4/lib/validator"
	"alterra-agmc-day-4/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	config.InitDB()
	database.CreateUser(&models.User{
		Email:    "user_1@email.com",
		Password: "Password",
	})
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
