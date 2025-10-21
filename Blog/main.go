package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := NewMySQL(dsn)
	if err != nil {
		log.Fatalf("initializing DB pool: %v", err)
	}
	AutoMigrate(db)

	r := gin.New()
	r.Use(RequestID(), Recover(), AccessLog())

	// Registers:
	//   POST /api/register
	//   POST /api/login
	//   GET  /api/me (protected, for quick verification)
	RegisterAuthRoutes(r, db, "/api")

	// POST /posts (auth required): Create a post with title and content.
	// GET /posts: List all posts.
	// GET /posts/:id: Get a single post by ID.
	// PUT /posts/:id (auth + author only): Update title/content.
	// DELETE /posts/:id (auth + author only): Delete post.
	RegisterPostRoutes(r, db, "/api")

	// POST base/posts/:id/comments (auth): Creates a comment for a post.
	// GET base/posts/:id/comments (public): Lists all comments for a post.
	RegisterCommentRoutes(r, db, "/api")

	// Optional health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
