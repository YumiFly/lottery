// routes/routes.go
package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
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

	// 用户相关路由
	r.POST("/login", controllers.Login)                                                                   //TODO 登录接口,需要改成根据地址登录，前端先连接钱包，后端验证
	r.POST("/customers/upload-photo", controllers.UploadPhoto)                                            // KYC 上传用户身份信息，上传用户头像等
	r.POST("/customers", middleware.ValidationMiddleware(&models.Customer{}), controllers.CreateCustomer) // KYC 用户注册接口
	r.GET("/customers", controllers.GetCustomers)                                                         // 获取所有用户，管理员员使用，需要分页，可以和下面的接口合并
	r.GET("/customers/:customer_address", controllers.GetCustomerByAddress)                               // 根据用户地址获取用户信息，需要验证用户身份
	r.GET("/customers/roles", controllers.GetRoleList)                                                    // 获取用户角色，需要验证用户身份

	r.POST("/lottery/types/v2", controllers.NewLotteryType)
	r.GET("/lottery/types/v2", controllers.ListLotteryTypes)

	r.POST("/lottery/lottery/v2", controllers.NewLottery)      // 创建彩票
	r.GET("/lottery/lottery/v2", controllers.ListAllLotteries) // 获取所有彩票信息

	r.POST("/lottery/issues/v2", controllers.NewLotteryIssue) // 发行彩票
	r.GET("/lottery/issues/v2", controllers.ListAllIssues)    // 通过分页获取所有发行信息

	r.POST("/lottery/tickets/v2", controllers.NewPurchaseTicket)  // 购买彩票
	r.GET("lottery/tickets/v2", controllers.ListPurchasedTickets) // 获取用户购买过的彩票信息

	r.POST("/lottery/draw/v2", controllers.NewDrawLottery) // 开奖

	r.GET("/lottery/winners/v2", controllers.ListWinners) // 获取近期得奖的用户信息

	r.GET("/lottery/pools/v2", controllers.CountIssuePools) // 获取彩票所有奖池总额

	// 静态文件服务（用于访问 uploads 目录下的文件）
	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/verify", controllers.VerifyCustomer)
	}
}
