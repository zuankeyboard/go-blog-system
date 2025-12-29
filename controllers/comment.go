package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论（需认证）
// @Summary 发表评论
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param post_id path int true "文章ID"
// @Param comment body struct{Content string} true "评论内容"
// @Success 200 {object} gin.H{"message":"评论成功","data":{}}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"未认证"}
// @Failure 404 {object} gin.H{"error":"文章不存在"}
// @Router /api/posts/{post_id}/comments [post]
func CreateComment(c *gin.Context) {
	// 1. 获取当前登录用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息，请先登录"})
		return
	}

	// 2. 从查询参数获取文章ID（替代URL路径参数）
	postIdStr := c.Query("post_id")
	if postIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文章ID（post_id）"})
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID格式错误"})
		return
	}

	// 3. 校验文章是否存在
	var post models.Post
	if err := config.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在，无法评论"})
		return
	}

	// 4. 绑定评论内容参数
	var req struct {
		Content string `json:"content" binding:"required,min=1"` // 评论内容非空，至少1个字符
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 5. 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userId.(uint),
		PostID:  uint(postId),
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败: " + err.Error()})
		return
	}

	// 6. 加载关联的用户+文章信息（关键修改）
	if err := config.DB.Preload("User").Preload("Post").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载评论关联信息失败: " + err.Error()})
		return
	}

	// 7. 转换时间为本地时区（新增文章时间转换）
	comment.CreatedAt = comment.CreatedAt.Local()
	comment.User.CreatedAt = comment.User.CreatedAt.Local()
	comment.User.UpdatedAt = comment.User.UpdatedAt.Local()
	comment.Post.CreatedAt = comment.Post.CreatedAt.Local()
	comment.Post.UpdatedAt = comment.Post.UpdatedAt.Local()

	// 8. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "评论发表成功",
		"data":    comment,
	})
}

// GetComments 获取某篇文章的所有评论
// @Summary 获取文章评论列表
// @Produce json
// @Param post_id path int true "文章ID"
// @Success 200 {object} gin.H{"data":[]models.Comment}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 404 {object} gin.H{"error":"文章不存在"}
// @Router /api/posts/{post_id}/comments [get]
func GetComments(c *gin.Context) {
	// 1. 从查询参数获取文章ID
	postIdStr := c.Query("post_id")
	if postIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文章ID（post_id）"})
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID格式错误"})
		return
	}

	// 2. 校验文章是否存在
	var post models.Post
	if err := config.DB.Where("id = ?", postId).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 3. 查询该文章的所有评论（关联加载评论用户+文章信息）
	var comments []models.Comment
	// 关键修改：Preload("User") + Preload("Post") 同时加载用户和文章关联
	if err := config.DB.Preload("User").Preload("Post").Where("post_id = ?", postId).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败: " + err.Error()})
		return
	}

	// 4. 转换所有评论的时间为本地时区（可选）
	for i := range comments {
		comments[i].CreatedAt = comments[i].CreatedAt.Local()
		comments[i].User.CreatedAt = comments[i].User.CreatedAt.Local()
		comments[i].User.UpdatedAt = comments[i].User.UpdatedAt.Local()
		// 新增：转换文章时间为本地时区
		comments[i].Post.CreatedAt = comments[i].Post.CreatedAt.Local()
		comments[i].Post.UpdatedAt = comments[i].Post.UpdatedAt.Local()
	}

	// 5. 返回评论列表
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}
