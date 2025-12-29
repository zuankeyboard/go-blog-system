package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"go-blog-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
// @Summary 用户注册
// @Accept json
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} gin.H{"message":"注册成功","data":{}}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 409 {object} gin.H{"error":"用户名/邮箱已存在"}
// @Router /api/register [post]
func Register(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查用户名/邮箱是否已存在
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
		return
	}

	// 创建用户（密码会通过BeforeCreate钩子自动加密）
	newUser := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		return
	}

	// 返回注册成功信息（隐藏密码）
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"user_id":  newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
		},
	})
}

// Login 用户登录接口
// @Summary 用户登录
// @Accept json
// @Produce json
// @Param user body struct{Username string;Password string} true "登录信息"
// @Success 200 {object} gin.H{"message":"登录成功","data":{"token":"xxx","user_id":1,"username":"xxx"}}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"用户名或密码错误"}
// @Router /api/login [post]
func Login(c *gin.Context) {
	// 绑定登录参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查询用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败: " + err.Error()})
		return
	}

	// 返回登录成功信息（包含Token）
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"user_id":  user.ID,
			"username": user.Username,
		},
	})
}
