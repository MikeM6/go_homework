package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 评论功能
// 实现评论的创建功能，已认证的用户可以对文章发表评论。
// 实现评论的读取功能，支持获取某篇文章的所有评论列表。

type CommentCreateRequest struct {
	PostID  uint64 `json:"postid" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// RegisterCommentRoutes 注册评论相关路由
//
// POST /api/posts/:id/comments  (需要认证)  创建评论
// GET  /api/posts/:id/comments  (公开)      获取某篇文章的所有评论
func RegisterCommentRoutes(r *gin.Engine, db *gorm.DB, base string) {
	g := r.Group(base)

	// 创建评论
	g.POST("/posts/comments", AuthMiddleware(), func(c *gin.Context) {
		var req CommentCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			BadRequest(c, err.Error())
			return
		}

		// 确认文章存在
		var post Post
		if err := db.First(&post, req.PostID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "post not found")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}

		uid, _ := c.Get("currentUserID")
		comment := &Comment{Content: req.Content, UserID: uid.(uint), PostID: post.ID}
		if err := db.Create(comment).Error; err != nil {
			c.Error(err)
			Internal(c, "failed to create comment")
			return
		}
		c.JSON(http.StatusCreated, comment)
	})

	// 获取某篇文章的评论列表
	g.GET("/posts/:id/comments", func(c *gin.Context) {
		postID, err := parseUintParam(c.Param("id"))
		if err != nil {
			c.Error(err)
			BadRequest(c, "invalid id")
			return
		}

		// 可选：确认文章存在，返回更明确的 404
		var post Post
		if err := db.Select("id").First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				NotFound(c, "post not found")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}

		var comments []Comment
		if err := db.Where("post_id = ?", postID).Order("id ASC").Find(&comments).Error; err != nil {
			c.Error(err)
			Internal(c, "database error")
			return
		}
		c.JSON(http.StatusOK, comments)
	})
}
