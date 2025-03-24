// routes/routes.go
package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.GlobalMiddleware())

	r.POST("/login", controllers.Login)
	// 将 customer 对象存入上下文，避免重复绑定
	r.POST("/customers", middleware.ValidationMiddleware(&models.Customer{}), controllers.CreateCustomer)
	r.GET("/customers", controllers.GetCustomers)
	r.GET("/customers/:customer_address", controllers.GetCustomerByAddress)

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/verify", controllers.VerifyCustomer)
	}
}
