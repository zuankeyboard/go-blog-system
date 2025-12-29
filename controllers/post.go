package controllers

import (
	"go-blog-system/config"
	"go-blog-system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建文章（需认证）
// @Summary 创建文章
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param post body struct{Title string;Content string} true "文章信息"
// @Success 200 {object} gin.H{"message":"创建成功","data":{}}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"未认证"}
// @Router /api/posts [post]
func CreatePost(c *gin.Context) {
	// 从上下文获取当前登录用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	// 绑定请求参数
	var req struct {
		Title   string `json:"title" binding:"required,min=1,max=200"` // 标题非空，长度1-200
		Content string `json:"content" binding:"required"`             // 内容非空
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId.(uint), // 关联当前登录用户
	}
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	// 关键修复：创建后，主动加载关联的 User 信息
	// Preload("User") 会根据 post.UserID 查询 users 表，填充 User 字段
	if err := config.DB.Preload("User").First(&post, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载作者信息失败: " + err.Error()})
		return
	}

	// 返回创建结果
	c.JSON(http.StatusOK, gin.H{
		"message": "文章创建成功",
		"data":    post,
	})
}

// GetPosts 获取所有文章列表
// @Summary 获取文章列表
// @Produce json
// @Success 200 {object} gin.H{"data":[]models.Post}
// @Router /api/posts [get]
func GetPosts(c *gin.Context) {
	// 定义文章列表切片
	var posts []models.Post
	// 查询所有文章，关联加载作者信息（Preload("User")）
	if err := config.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败: " + err.Error()})
		return
	}

	// 返回文章列表
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

// GetPost 获取单篇文章详情
// @Summary 获取单篇文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} gin.H{"data":models.Post}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 404 {object} gin.H{"error":"文章不存在"}
// @Router /api/posts/{id} [get]
func GetPost(c *gin.Context) {
	// 解析URL中的文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID格式错误"})
		return
	}

	// 查询文章（关联作者信息）
	var post models.Post
	if err := config.DB.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 返回文章详情
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

// UpdatePost 更新文章（仅作者可操作）
// @Summary 更新文章
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "文章ID"
// @Param post body struct{Title string;Content string} true "更新信息"
// @Success 200 {object} gin.H{"message":"更新成功","data":{}}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"未认证"}
// @Failure 403 {object} gin.H{"error":"无权限操作"}
// @Failure 404 {object} gin.H{"error":"文章不存在"}
// @Router /api/posts/{id} [put]
func UpdatePost(c *gin.Context) {
	// 获取当前登录用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	// 解析文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID格式错误"})
		return
	}

	// 查询文章是否存在，且作者是当前用户
	var post models.Post
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&post).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新该文章（文章不存在或非本人创建）"})
		return
	}

	// 绑定更新参数
	var req struct {
		Title   string `json:"title" binding:"omitempty,min=1,max=200"` // 可选，更新时非空
		Content string `json:"content" binding:"omitempty"`             // 可选
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 更新文章字段（仅更新非空的字段）
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	// 保存更新
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败: " + err.Error()})
		return
	}

	// 返回更新结果
	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"data":    post,
	})
}

// DeletePost 删除文章（仅作者可操作）
// @Summary 删除文章
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "文章ID"
// @Success 200 {object} gin.H{"message":"删除成功"}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"未认证"}
// @Failure 403 {object} gin.H{"error":"无权限操作"}
// @Failure 404 {object} gin.H{"error":"文章不存在"}
// @Router /api/posts/{id} [delete]
func DeletePost(c *gin.Context) {
	// 获取当前登录用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}

	// 解析文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID格式错误"})
		return
	}

	// 查询文章是否存在，且作者是当前用户
	var post models.Post
	if err := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&post).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除该文章（文章不存在或非本人创建）"})
		return
	}

	// 删除文章（GORM默认软删除，如需物理删除可使用 Unscoped().Delete()）
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败: " + err.Error()})
		return
	}

	// 返回删除结果
	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
	})
}
