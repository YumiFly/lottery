// tests/auth_test.go
package tests

import (
	"backend/controllers"
	"backend/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
	suite := SetupTestDB()
	defer suite.TearDown()

	t.Run("Login", func(t *testing.T) {
		result, err := services.Login("0xTestAddress123", "127.0.0.1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "0xTestAddress123", result.Customer.CustomerAddress)
		assert.NotEmpty(t, result.Token)
	})

	t.Run("LoginNonExistent", func(t *testing.T) {
		result, err := services.Login("0xNewAddress789", "127.0.0.1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "0xNewAddress789", result.Customer.CustomerAddress)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		newToken, err := services.RefreshToken("test@example.com", "lottery_admin")
		assert.NoError(t, err)
		assert.NotEmpty(t, newToken)
	})
}

func TestAuthController(t *testing.T) {
	suite := SetupTestDB()
	defer suite.TearDown()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth/login", controllers.Login)

	t.Run("Login", func(t *testing.T) {
		body := map[string]string{
			"wallet_address": "0xTestAddress123",
		}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "test@example.com")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "0xTestAddress123")
	})

	t.Run("LoginInvalid", func(t *testing.T) {
		body := map[string]string{}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "invalid@example.com")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
