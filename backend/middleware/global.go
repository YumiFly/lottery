// middleware/global.go
package middleware

import (
	"net/http"
	"time"

	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// GlobalMiddleware 全局中间件，使用 logrus 记录日志
func GlobalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录请求日志
		duration := time.Since(start)
		status := c.Writer.Status()

		utils.Logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   status,
			"duration": duration.String(),
			"ip":       c.ClientIP(),
		}).Info("Request processed")

		// 如果有错误，记录详细日志
		if status >= 400 {
			if errs, exists := c.Get("errors"); exists {
				utils.Logger.WithFields(logrus.Fields{
					"errors": errs,
				}).Error("Request failed")
			}
		}
	}
}

// ValidationMiddleware 校验中间件，记录校验错误
func ValidationMiddleware(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if obj == nil {
			c.Next()
			return
		}

		// 绑定 JSON 数据
		if err := c.ShouldBindJSON(obj); err != nil {
			utils.Logger.WithField("error", err.Error()).Warn("Invalid JSON format")
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid JSON format", err.Error()))
			c.Abort()
			return
		}

		// 验证结构体
		validate := utils.NewValidator()
		if err := validate.Struct(obj); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, err.Field()+" "+err.Tag()+" validation failed")
			}
			utils.Logger.WithField("errors", errors).Warn("Validation failed")
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeValidationFailed, "Validation failed", errors))
			c.Abort()
			return
		}

		// 将验证后的对象存入上下文
		c.Set("validated_obj", obj)

		c.Next()
	}
}
