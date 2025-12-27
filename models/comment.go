package models

import "gorm.io/gorm"

// Comment 对应 comments 表，存储文章评论信息
type Comment struct {
	gorm.Model        // 内置字段：ID、CreatedAt、UpdatedAt、DeletedAt
	Content    string `gorm:"type:text;not null" json:"content"` // 评论内容，非空
	UserID     uint   `gorm:"not null" json:"user_id"`           // 关联评论用户ID（外键）
	PostID     uint   `gorm:"not null" json:"post_id"`           // 关联文章ID（外键）
	// 关联模型：查询时可加载评论用户/所属文章信息
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}
