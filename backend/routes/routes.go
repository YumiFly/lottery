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

	// 静态文件服务（用于访问 uploads 目录下的文件）
	r.Static("/uploads", "./uploads")

	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/verify", controllers.VerifyCustomer)
	}
}
