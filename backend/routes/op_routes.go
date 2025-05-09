// routes/routes.go
package routes

import (
	"backend/controllers"
	"backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupOpRoutes(r *gin.Engine) {
	r.Use(middleware.GlobalMiddleware())

	// 配置 CORS 中间件
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // 允许的来源，设置为你的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 暴露给前端的响应头
		AllowCredentials: true,                                                // 是否允许发送凭证（如 cookies）
		MaxAge:           12 * time.Hour,                                      // 预检请求（OPTIONS）的缓存时间
	}

	// 应用 CORS 中间件
	r.Use(cors.New(config))

	r.POST("/login", controllers.Login)

	// 稳定币管理相关
	// 增加/设置稳定币
	r.POST("/stablecoin", controllers.SetStableCoin)
	// 删除稳定币
	r.DELETE("/stablecoin", controllers.RemoveStableCoin)

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/verify", controllers.VerifyCustomer)
	}
}
