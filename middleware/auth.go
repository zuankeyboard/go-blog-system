package middleware

import (
	"go-blog-system/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件（接收JWT密钥参数）
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			utils.Log.Warnf("未携带Token, ip: %s", c.ClientIP())
			utils.Unauthorized(c, "未携带Token，请先登录")
			return
		}

		// 去掉Bearer前缀
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		// 解析Token
		claims, err := utils.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			utils.Log.Warnf("Token解析失败: %v, ip: %s", err, c.ClientIP())
			utils.Unauthorized(c, "Token无效或已过期")
			return
		}

		// 设置上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
