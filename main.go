package main

import (
	"go-blog-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 连接 SQLite 数据库（文件名为 blog.db，不存在则自动创建）
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		panic("表结构迁移失败：" + err.Error())
	}
}
