// tests/user_test.go
package tests

import (
	"backend/controllers"
	"backend/models"
	"backend/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	suite := SetupTestDB()
	defer suite.TearDown()

	t.Run("GetCustomers", func(t *testing.T) {
		customers, err := services.GetCustomers()
		assert.NoError(t, err)
		assert.Len(t, customers, 1)
		assert.Equal(t, "0xTestAddress123", customers[0].CustomerAddress)
	})

	t.Run("GetCustomerByAddress", func(t *testing.T) {
		customer, err := services.GetCustomerByAddress("0xTestAddress123")
		assert.NoError(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, "0xTestAddress123", customer.CustomerAddress)
		assert.Equal(t, "test@example.com", customer.KYCData.Email)

		// 测试不存在的用户
		_, err = services.GetCustomerByAddress("0xNonExistent")
		assert.Error(t, err)
	})

	t.Run("CreateCustomer", func(t *testing.T) {
		newCustomer := &models.Customer{
			CustomerAddress: "0xNewAddress456",
			RoleID:          1,
			KYCData: models.KYCData{
				CustomerAddress: "0xNewAddress456",
				Email:           "new@example.com",
				Name:            "New User",
			},
		}
		err := services.CreateCustomer(newCustomer)
		assert.NoError(t, err)

		// 验证创建结果
		customer, err := services.GetCustomerByAddress("0xNewAddress456")
		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", customer.KYCData.Email)
	})
}

func TestUserController(t *testing.T) {
	suite := SetupTestDB()
	defer suite.TearDown()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/customers", controllers.GetCustomers)
	r.GET("/customers/:customer_address", controllers.GetCustomerByAddress)

	t.Run("GetCustomers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/customers", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "0xTestAddress123")
	})

	t.Run("GetCustomerByAddress", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/customers/0xTestAddress123", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test@example.com")

		// 测试不存在的用户
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/customers/0xNonExistent", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
