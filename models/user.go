package models

import (
	"gorm.io/gorm"
)

// User 对应 users 表，存储用户核心信息
type User struct {
	// GORM 内置字段：ID（主键）、CreatedAt、UpdatedAt、DeletedAt（软删除）
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名，唯一且非空
	Password string `gorm:"size:100;not null" json:"-"`                   // 密码（加密存储，前端不返回）
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`            // 邮箱，唯一
}

// // BeforeCreate GORM 钩子：创建用户前自动加密密码
// func (u *User) BeforeCreate(tx *gorm.DB) error {
// 	// 密码加密：使用 bcrypt 生成哈希值
// 	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = string(hashedPwd)
// 	return nil
// }

// // CheckPassword 验证密码是否正确
// func (u *User) CheckPassword(password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
// 	return err == nil
// }
