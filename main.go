package main

import (
	"go-blog-system/config"
	"go-blog-system/controllers"
	"go-blog-system/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 创建Gin引擎
	r := gin.Default()

	// 公开路由（无需认证）
	publicGroup := r.Group("/api")
	{
		publicGroup.POST("/register", controllers.Register) // 注册
		publicGroup.POST("/login", controllers.Login)       // 登录
	}

	// 私有路由（需要JWT认证）
	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware()) // 应用JWT中间件
	{
		// 后续添加需要认证的接口（如文章CRUD）
		privateGroup.GET("/profile", func(c *gin.Context) {
			// 从上下文获取用户信息
			userID := c.GetInt("user_id")
			username := c.GetString("username")
			c.JSON(http.StatusOK, gin.H{
				"message": "获取个人信息成功",
				"data": gin.H{
					"user_id":  userID,
					"username": username,
				},
			})
		})
	}

	// 启动服务（监听8080端口）
	if err := r.Run(":8080"); err != nil {
		panic("服务启动失败: " + err.Error())
	}
}
