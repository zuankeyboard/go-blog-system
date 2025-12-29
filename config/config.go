package config

import (
	"go-blog-system/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AppConfig 全局配置结构体
type AppConfig struct {
	JWTSecretKey string // JWT密钥
	TokenExpire  int    // Token过期时间（小时）
	DBFile       string // SQLite文件路径
}

// 全局DB实例
var DB *gorm.DB

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *AppConfig {
	return &AppConfig{
		JWTSecretKey: "blog-jwt-secret-2025",
		TokenExpire:  72,
		DBFile:       "blog.db",
	}
}

// InitDB 初始化数据库
func InitDB(cfg *AppConfig) {
	gormLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBFile), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatalf("[Config] 数据库连接失败: %v", err)
		panic("数据库连接失败: " + err.Error())
	}

	// 自动迁移表
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Printf("[Config] 表迁移失败: %v", err)
		panic("表迁移失败: " + err.Error())
	}
	log.Println("[Config] 数据库初始化成功")
}
