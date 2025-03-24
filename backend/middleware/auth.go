// middleware/auth.go
package middleware

import (
	"backend/config"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 定义验证 JWT 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Authorization header required", nil))
			c.Abort()
			return
		}
		// 提取 Bearer 令牌
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Invalid Authorization header", nil))
			c.Abort()
			return
		}
		token := authHeader[7:]
		// 验证 JWT 令牌
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Invalid token", err.Error()))
			c.Abort()
			return
		}
		// 将用户信息存入上下文
		c.Set("customer_address", claims["customer_address"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
