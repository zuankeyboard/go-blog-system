package config

import (
	"go-blog-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 全局DB实例
var DB *gorm.DB

// JWT配置常量
const (
	JWTSecretKey = "blog-jwt-secret-2025" // JWT签名密钥（可自定义）
	TokenExpire  = 72                     // Token过期时间（小时）
)

// 初始化数据库连接（SQLite）
func InitDB() {
	var err error
	// 连接SQLite数据库（文件为blog.db，不存在自动创建）
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移表结构（创建users/posts/comments表）
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		panic("表结构迁移失败: " + err.Error())
	}
}
