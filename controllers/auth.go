package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"go-blog-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}

	// 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log.Warnf("注册参数错误: %v, ip: %s", err, c.ClientIP())
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否存在
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		utils.Log.Warnf("用户名已存在: %s, ip: %s", req.Username, c.ClientIP())
		utils.Forbidden(c, "用户名已存在")
		return
	}

	// 检查邮箱是否存在
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		utils.Log.Warnf("邮箱已存在: %s, ip: %s", req.Email, c.ClientIP())
		utils.Forbidden(c, "邮箱已存在")
		return
	}

	// 创建用户
	newUser := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := config.DB.Create(&newUser).Error; err != nil {
		utils.Log.Errorf("创建用户失败: %v, ip: %s", err, c.ClientIP())
		utils.InternalError(c, "注册失败: "+err.Error())
		return
	}

	utils.Log.Infof("用户注册成功: %s, id: %d", req.Username, newUser.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"user_id":  newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
		},
	})
}

// Login 用户登录（接收JWT配置参数）
func Login(c *gin.Context, jwtSecret string, tokenExpire int) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log.Warnf("登录参数错误: %v, ip: %s", err, c.ClientIP())
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 查询用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Log.Warnf("用户不存在: %s, ip: %s", req.Username, c.ClientIP())
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.Log.Warnf("密码错误: %s, ip: %s", req.Username, c.ClientIP())
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, jwtSecret, tokenExpire)
	if err != nil {
		utils.Log.Errorf("生成Token失败: %v, user_id: %d", err, user.ID)
		utils.InternalError(c, "登录失败: "+err.Error())
		return
	}

	utils.Log.Infof("用户登录成功: %s, id: %d", req.Username, user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"user_id":  user.ID,
			"username": user.Username,
		},
	})
}
