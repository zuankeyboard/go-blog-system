package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"go-blog-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	// 获取当前用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required,min=1,max=100"`
		Content string `json:"content" binding:"required,min=1"`
	}

	// 绑定参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log.Warnf("创建文章参数错误: %v, user_id: %d", err, userId)
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId.(uint),
	}
	if err := config.DB.Create(&post).Error; err != nil {
		utils.Log.Errorf("创建文章失败: %v, user_id: %d", err, userId)
		utils.InternalError(c, "创建文章失败: "+err.Error())
		return
	}

	// 加载作者信息
	if err := config.DB.Preload("User").First(&post, post.ID).Error; err != nil {
		utils.Log.Warnf("加载文章作者信息失败: %v, post_id: %d", err, post.ID)
	}

	utils.Log.Infof("文章创建成功: post_id: %d, user_id: %d", post.ID, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "文章创建成功",
		"data":    post,
	})
}

// GetPosts 获取所有文章
func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		utils.Log.Errorf("获取文章列表失败: %v", err)
		utils.InternalError(c, "获取文章列表失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

// GetPost 获取单篇文章
func GetPost(c *gin.Context) {
	// 解析文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Log.Warnf("文章ID格式错误: %v, ip: %s", err, c.ClientIP())
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 查询文章
	var post models.Post
	if err := config.DB.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		utils.Log.Infof("文章不存在: id=%d, ip: %s", id, c.ClientIP())
		utils.NotFound(c, "文章不存在")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	// 获取当前用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	// 解析文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Log.Warnf("文章ID格式错误: %v, user_id: %d", err, userId)
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 查询文章并验证归属
	var post models.Post
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&post).Error; err != nil {
		utils.Log.Warnf("文章不存在或无权限: id=%d, user_id: %d", id, userId)
		utils.NotFound(c, "文章不存在或无修改权限")
		return
	}

	// 绑定更新参数
	var req struct {
		Title   string `json:"title" binding:"omitempty,min=1,max=100"`
		Content string `json:"content" binding:"omitempty,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log.Warnf("更新文章参数错误: %v, post_id: %d", err, id)
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 更新文章
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if err := config.DB.Save(&post).Error; err != nil {
		utils.Log.Errorf("更新文章失败: %v, post_id: %d", err, id)
		utils.InternalError(c, "更新文章失败: "+err.Error())
		return
	}

	// 重新加载作者信息
	config.DB.Preload("User").First(&post, post.ID)

	utils.Log.Infof("文章更新成功: post_id: %d, user_id: %d", id, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"data":    post,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	// 获取当前用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	// 解析文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Log.Warnf("文章ID格式错误: %v, user_id: %d", err, userId)
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 查询文章并验证归属
	var post models.Post
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&post).Error; err != nil {
		utils.Log.Warnf("文章不存在或无权限: id=%d, user_id: %d", id, userId)
		utils.NotFound(c, "文章不存在或无删除权限")
		return
	}

	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.Log.Errorf("删除文章失败: %v, post_id: %d", err, id)
		utils.InternalError(c, "删除文章失败: "+err.Error())
		return
	}

	utils.Log.Infof("文章删除成功: post_id: %d, user_id: %d", id, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
	})
}
