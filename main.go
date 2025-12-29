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
		// 文章公开读取接口
		publicGroup.GET("/posts", controllers.GetPosts)    // 获取所有文章
		publicGroup.GET("/posts/:id", controllers.GetPost) // 获取单篇文章
	}

	// 私有路由（需要JWT认证）
	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware()) // 应用JWT中间件
	{
		// 个人信息
		privateGroup.GET("/profile", func(c *gin.Context) {
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

		// 文章管理接口（需认证+权限）
		privateGroup.POST("/posts", controllers.CreatePost)       // 创建文章
		privateGroup.PUT("/posts/:id", controllers.UpdatePost)    // 更新文章
		privateGroup.DELETE("/posts/:id", controllers.DeletePost) // 删除文章
	}

	// 启动服务（监听8080端口）
	if err := r.Run(":8080"); err != nil {
		panic("服务启动失败: " + err.Error())
	}
}
