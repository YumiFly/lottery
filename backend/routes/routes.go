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

	r.POST("/login", controllers.Login)
	// 文件上传接口
	r.POST("/customers/upload-photo", controllers.UploadPhoto)
	r.POST("/customers", middleware.ValidationMiddleware(&models.Customer{}), controllers.CreateCustomer)
	r.GET("/customers", controllers.GetCustomers)
	r.GET("/customers/:customer_address", controllers.GetCustomerByAddress)

	// 彩票相关路由
	r.POST("/lottery/types", controllers.CreateLotteryType) // 创建彩票类型
	r.POST("/lottery/lottery", controllers.CreateLottery)   // 创建彩票
	r.POST("/lottery/issues", controllers.CreateIssue)      // 发行彩票
	r.POST("/lottery/tickets", controllers.PurchaseTicket)  // 购买彩票
	r.POST("/lottery/draw", controllers.DrawLottery)        // 开奖

	// 获取彩票类型
	r.GET("/lottery/types", controllers.GetAllLotteryTypes) // 获取所有彩票类型
	// 获取彩票信息
	r.GET("/lottery/lottery", controllers.GetAllLottery) // 获取所有彩票信息
	// 根据彩票类型获取彩票信息
	r.GET("/lottery/lottery/:lottery_type", controllers.GetLotteryByType) // 根据彩票类型获取彩票信息

	// 获取即将开奖的彩票信息
	r.GET("/lottery/upcoming-issues", controllers.GetUpcomingIssues)

	// 获取根据彩票ID获取最近的发行信息
	r.GET("/lottery/issues/latest/:lottery_id", controllers.GetLatestIssueByLotteryID)

	// 获取彩票所有奖池总额
	r.GET("/lottery/pools", controllers.GetAllPools)

	// 获取用户购买过的彩票信息
	r.GET("lottery/tickets/customer/:customer_address", controllers.GetPurchasedTicketsByCustomerAddress)

	// 获取开奖信息
	r.GET("/lottery/draw/:issue_id", controllers.GetDrawnLotteryByIssueID)
	// 获取期号信息
	r.GET("/lottery/issues/:issue_id", controllers.GetIssueByID)

	// 获取近期开奖彩票信息和开奖结果
	r.GET("/lottery/draw/latest", controllers.GetLatestDrawnLottery)

	// 获取近期得奖的用户信息
	r.GET("/lottery/recent-winners", controllers.GetRecentWinners)

	// 静态文件服务（用于访问 uploads 目录下的文件）
	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/verify", controllers.VerifyCustomer)
	}
}
