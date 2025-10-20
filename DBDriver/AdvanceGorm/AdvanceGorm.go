// Package advancegorm
package advancegorm

import (
	"database/sql"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。

// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。

// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

type User struct {
	gorm.Model
	Name  string `gorm:"size:64;not null"`
	Email string `gorm:"size:128;uniqueIndex;not null"`

	// has many
	Posts []Post

	PostCount int `gorm:"default:0"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null"`
	Content string `gorm:"type:text"`

	// Foreign key to User
	UserID uint `gorm:"index;not null"`

	// One-to-many: Post has many Comments
	Comments      []Comment
	CommentStatus string `gorm:"size:32;default:''"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`

	// Foreign key to Post
	PostID uint `gorm:"index;not null"`

	// (Optional) If you want to track commenter user, keep UserID
	// Remove if not needed by your schema requirements
	UserID uint `gorm:"index"`
}

// DeleteOneComment : delete a single comment for post_id=2
// Ensures AfterDelete receives the correct PostID.
func DeleteOneComment(db *gorm.DB) error {
	var c Comment
	if err := db.Where("post_id = ?", 2).Order("id ASC").First(&c).Error; err != nil {
		return err // no comment found, or DB error
	}
	return db.Delete(&c).Error
}

// CreateOnePost : Test AfterCreate Hook
func CreateOnePost(db *gorm.DB) error {
	p2 := Post{Title: "Test AfterCreate Hook", Content: "Test AfterCreate Hook", UserID: 2}
	if err := db.Where(Post{Title: p2.Title, UserID: 2}).FirstOrCreate(&p2).Error; err != nil {
		return err
	}
	return nil
}

// AfterCreate hook for Post:
// When a post is created, automatically increment the author's PostCount.
func (p *Post) AfterCreate(tx *gorm.DB) error {
	if p == nil || p.UserID == 0 {
		return nil
	}
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		// atomic increment, avoid triggering other hooks
		UpdateColumn("post_count", gorm.Expr("COALESCE(post_count, 0) + 1")).
		Error
}

// AfterDelete hook for Comment:
// When a comment is deleted, if the post has zero remaining comments,
// set the post's CommentStatus to "无评论".
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	if c == nil || c.PostID == 0 {
		return nil
	}
	var cnt int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error
	}
	return nil
}

func NewMySQL(dsn string) (*gorm.DB, error) {
	gLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gLogger})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

func BatchInsertBlogData(db *gorm.DB) error {
	tx := db.Begin()

	// Users
	alice := User{Name: "Alice", Email: "alice@example.com"}
	if err := tx.Where(User{Email: alice.Email}).FirstOrCreate(&alice).Error; err != nil {
		tx.Rollback()
		return err
	}
	bob := User{Name: "Bob", Email: "bob@example.com"}
	if err := tx.Where(User{Email: bob.Email}).FirstOrCreate(&bob).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Alice's posts
	p1 := Post{Title: "Go with GORM", Content: "Intro to GORM and associations", UserID: alice.ID}
	if err := tx.Where(Post{Title: p1.Title, UserID: alice.ID}).FirstOrCreate(&p1).Error; err != nil {
		tx.Rollback()
		return err
	}
	p2 := Post{Title: "Concurrency Patterns", Content: "Goroutines and channels", UserID: alice.ID}
	if err := tx.Where(Post{Title: p2.Title, UserID: alice.ID}).FirstOrCreate(&p2).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Bob's post
	p3 := Post{Title: "SQL Tips", Content: "Indexes and query plans", UserID: bob.ID}
	if err := tx.Where(Post{Title: p3.Title, UserID: bob.ID}).FirstOrCreate(&p3).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Comments for p1 (3 comments → p1 should be “most commented”)
	for _, content := range []string{"Great intro!", "Helped me a lot.", "More examples please."} {
		c := Comment{Content: content, PostID: p1.ID}
		if err := tx.Where(Comment{Content: c.Content, PostID: p1.ID}).FirstOrCreate(&c).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Comments for p2 (1)
	if err := tx.Where(Comment{Content: "Nice patterns.", PostID: p2.ID}).FirstOrCreate(&Comment{Content: "Nice patterns.", PostID: p2.ID}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Comments for p3 (2)
	for _, content := range []string{"Indexing is key.", "Use EXPLAIN often."} {
		c := Comment{Content: content, PostID: p3.ID}
		if err := tx.Where(Comment{Content: c.Content, PostID: p3.ID}).FirstOrCreate(&c).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error

}

// StdDB exposes the underlying *sql.DB for lifecycle control (e.g., Close).
func StdDB(db *gorm.DB) (*sql.DB, error) { return db.DB() }

func GetUserPostsWithComments(db *gorm.DB, userID uint) ([]Post, error) {
	var posts []Post
	err := db.
		Where("user_id = ?", userID).
		Preload("Comments").
		Order("id ASC").
		Find(&posts).Error
	return posts, err
}

func GetPostWithMostComments(db *gorm.DB) (Post, int64, error) {
	// First compute the post id with max comment count
	var row struct {
		ID           uint  `gorm:"column:id"`
		CommentCount int64 `gorm:"column:comment_count"`
	}

	err := db.Table("posts").
		Select("posts.id, COUNT(comments.id) AS comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return Post{}, 0, err
	}
	// If no rows (no posts), return zero values
	if row.ID == 0 {
		return Post{}, 0, nil
	}
	// Fetch the full post record and preload comments
	var post Post
	if err := db.Preload("Comments").First(&post, row.ID).Error; err != nil {
		return Post{}, 0, err
	}
	return post, row.CommentCount, nil
}
