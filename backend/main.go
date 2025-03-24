// main.go
package main

import (
	"backend/config"
	"backend/db"
	"backend/routes"
	"backend/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Lottery Backend API
// @version 1.0
// @description API for lottery backend with KYC management
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()
	db.InitDB()
	utils.InitLogger()
	utils.InitCache()

	r := gin.Default()
	routes.SetupRoutes(r)

	// 服务 Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
