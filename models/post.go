package models

import "gorm.io/gorm"

// Post 对应 posts 表，存储博客文章信息
type Post struct {
	gorm.Model        // 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt
	Title      string `gorm:"size:200;not null" json:"title"`    // 文章标题，非空
	Content    string `gorm:"type:text;not null" json:"content"` // 文章内容，文本类型
	UserID     uint   `gorm:"not null" json:"user_id"`           // 关联用户ID（外键）
	// 关联 User 模型（一对一），查询时可通过 Preload("User") 加载用户信息
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
