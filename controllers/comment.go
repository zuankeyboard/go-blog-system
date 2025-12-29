package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"go-blog-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 获取当前用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未获取到用户信息")
		return
	}

	// 解析文章ID（查询参数）
	postIdStr := c.Query("post_id")
	if postIdStr == "" {
		utils.BadRequest(c, "缺少文章ID（post_id）")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		utils.Log.Warnf("文章ID格式错误: %v, user_id: %d", err, userId)
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 校验文章是否存在
	var post models.Post
	if err := config.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		utils.Log.Warnf("文章不存在: id=%d, user_id: %d", postId, userId)
		utils.NotFound(c, "文章不存在，无法评论")
		return
	}

	// 绑定评论内容
	var req struct {
		Content string `json:"content" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Log.Warnf("评论参数错误: %v, user_id: %d", err, userId)
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userId.(uint),
		PostID:  uint(postId),
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		utils.Log.Errorf("创建评论失败: %v, user_id: %d", err, userId)
		utils.InternalError(c, "发表评论失败: "+err.Error())
		return
	}

	// 加载关联信息
	if err := config.DB.Preload("User").Preload("Post").First(&comment, comment.ID).Error; err != nil {
		utils.Log.Warnf("加载评论关联信息失败: %v, comment_id: %d", err, comment.ID)
	}

	utils.Log.Infof("评论创建成功: comment_id: %d, post_id: %d, user_id: %d", comment.ID, postId, userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "评论发表成功",
		"data":    comment,
	})
}

// GetComments 获取文章评论列表
func GetComments(c *gin.Context) {
	// 解析文章ID（查询参数）
	postIdStr := c.Query("post_id")
	if postIdStr == "" {
		utils.BadRequest(c, "缺少文章ID（post_id）")
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		utils.Log.Warnf("文章ID格式错误: %v, ip: %s", err, c.ClientIP())
		utils.BadRequest(c, "文章ID格式错误")
		return
	}

	// 校验文章是否存在
	var post models.Post
	if err := config.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		utils.Log.Warnf("文章不存在: id=%d, ip: %s", postId, c.ClientIP())
		utils.NotFound(c, "文章不存在")
		return
	}

	// 查询评论
	var comments []models.Comment
	if err := config.DB.Preload("User").Preload("Post").Where("post_id = ?", postId).Order("created_at DESC").Find(&comments).Error; err != nil {
		utils.Log.Errorf("获取评论列表失败: %v, post_id: %d", err, postId)
		utils.InternalError(c, "获取评论列表失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}
