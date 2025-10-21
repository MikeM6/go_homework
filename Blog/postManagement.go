package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 文章管理功能
// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
// 实现文章的删除功能，只有文章的作者才能删除自己的文章。

type PostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required"`
}

type PostUpdateRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// RegisterPostRoutes 注册文章相关路由
func RegisterPostRoutes(r *gin.Engine, db *gorm.DB, base string) {
	g := r.Group(base)

	// 创建文章（需要认证）
	g.POST("/posts", AuthMiddleware(), func(c *gin.Context) {
		var req PostCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			BadRequest(c, err.Error())
			return
		}
		uid, _ := c.Get("currentUserID")
		post := &Post{Title: req.Title, Content: req.Content, UserID: uid.(uint)}
		if err := db.Create(post).Error; err != nil {
			c.Error(err)
			Internal(c, "database error: create post fail")
			return
		}
		c.JSON(http.StatusCreated, post)
	})

	// 获取文章列表（公开）
	g.GET("/posts", func(c *gin.Context) {
		var posts []Post
		if err := db.Order("id DESC").Find(&posts).Error; err != nil {
			Internal(c, "database error: get post fail")
			return
		}
		c.JSON(http.StatusOK, posts)
	})

	// 获取文章详情（公开）
	g.GET("/posts/:id", func(c *gin.Context) {
		id, err := parseUintParam(c.Param("id"))
		if err != nil {
			c.Error(err)
			BadRequest(c, "invalid id")
			return
		}
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "post not found")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}
		c.JSON(http.StatusOK, post)
	})

	// 更新文章（仅作者）
	g.PUT("/posts/:id", AuthMiddleware(), func(c *gin.Context) {
		id, err := parseUintParam(c.Param("id"))
		if err != nil {
			c.Error(err)
			BadRequest(c, "invalid id")
			return
		}
		var req PostUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			BadRequest(c, err.Error())
			return
		}
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "post not found")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}
		uid, _ := c.Get("currentUserID")
		if post.UserID != uid.(uint) {
			Forbidden(c, "permission denied")
			return
		}
		updates := map[string]interface{}{}
		if req.Title != nil {
			updates["title"] = *req.Title
		}
		if req.Content != nil {
			updates["content"] = *req.Content
		}
		if len(updates) == 0 {
			WellRequest(c, "no fields to update")
			return
		}
		if err := db.Model(&post).Updates(updates).Error; err != nil {
			c.Error(err)
			Internal(c, "failed to update post")
			return
		}
		c.JSON(http.StatusOK, post)
	})

	// 删除文章（仅作者）
	g.DELETE("/posts/:id", AuthMiddleware(), func(c *gin.Context) {
		id, err := parseUintParam(c.Param("id"))
		if err != nil {
			c.Error(err)
			BadRequest(c, "invalid id")
			return
		}
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "post not found")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}
		uid, _ := c.Get("currentUserID")
		if post.UserID != uid.(uint) {
			Forbidden(c, "permission denied")
			return
		}
		if err := db.Delete(&post).Error; err != nil {
			c.Error(err)
			Internal(c, "failed to delete post")
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	})
}

func parseUintParam(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return uint(v), err
}
