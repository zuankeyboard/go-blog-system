package middleware

import (
	"go-blog-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
// 验证Token有效性，通过后将用户信息存入上下文
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token（格式：Bearer <token> 或直接 <token>）
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Token"})
			c.Abort() // 终止请求链
			return
		}

		// 去除Bearer前缀（如果有）
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效或已过期: " + err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续接口使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 继续执行后续中间件/接口
		c.Next()
	}
}
