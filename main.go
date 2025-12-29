package main

import (
	"go-blog-system/config"
	"go-blog-system/controllers"
	"go-blog-system/middleware"
	"go-blog-system/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化日志
	utils.InitLogger()
	utils.Log.Info("博客系统启动中...")

	// 2. 初始化配置
	appCfg := config.NewDefaultConfig()

	// 3. 初始化数据库
	config.InitDB(appCfg)

	// 4. Gin引擎配置
	r := gin.New()
	r.Use(utils.GinLogger()) // 自定义日志中间件
	r.Use(gin.Recovery())    // 异常恢复

	// 跨域中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 5. 路由配置
	publicGroup := r.Group("/api")
	{
		// 用户接口
		publicGroup.POST("/register", controllers.Register)
		publicGroup.POST("/login", func(c *gin.Context) {
			controllers.Login(c, appCfg.JWTSecretKey, appCfg.TokenExpire)
		})

		// 文章接口
		publicGroup.GET("/posts", controllers.GetPosts)
		publicGroup.GET("/posts/:id", controllers.GetPost)

		// 评论接口
		publicGroup.GET("/comments", controllers.GetComments)
	}

	// 私有路由（需要JWT认证）
	privateGroup := r.Group("/api")
	privateGroup.Use(middleware.JWTAuthMiddleware(appCfg.JWTSecretKey))
	{
		// 个人信息
		privateGroup.GET("/profile", func(c *gin.Context) {
			userID := c.GetInt("user_id")
			username := c.GetString("username")
			c.JSON(200, gin.H{
				"message": "获取个人信息成功",
				"data": gin.H{
					"user_id":  userID,
					"username": username,
				},
			})
		})

		// 文章接口
		privateGroup.POST("/posts", controllers.CreatePost)
		privateGroup.PUT("/posts/:id", controllers.UpdatePost)
		privateGroup.DELETE("/posts/:id", controllers.DeletePost)

		// 评论接口
		privateGroup.POST("/comments", controllers.CreateComment)
	}

	// 6. 启动服务
	utils.Log.Info("博客系统启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		utils.Log.Fatalf("服务启动失败: %v", err)
		panic("服务启动失败: " + err.Error())
	}
}
