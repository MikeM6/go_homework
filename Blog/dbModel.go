package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 2. 数据库设计与模型定义
// 设计数据库表结构，至少包含以下几个表：
// users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
// posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
// comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
// 使用 GORM 定义对应的 Go 模型结构体。

// User represents the users table.
type User struct {
	gorm.Model
	Username string `gorm:"size:64;uniqueIndex;not null"`
	Password string `gorm:"size:256;not null"`
	Email    string `gorm:"size:128;uniqueIndex;not null"`

	// Relations
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

// Post represents the posts table.
type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null"`
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null;index"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Relations
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment represents the comments table.
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null;index"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID  uint   `gorm:"not null;index"`
	Post    Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName (explicit for clarity; GORM would infer these by default)
func (User) TableName() string    { return "users" }
func (Post) TableName() string    { return "posts" }
func (Comment) TableName() string { return "comments" }

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

func NewMySQL(dsn string) (*gorm.DB, error) {
	// set logger level, output info level log.
	gLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gLogger,
	})
	if err != nil {
		return nil, err
	}

	// db.DB() gives access to that low-level object so you can tune the connection pool.
	// set connect pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func Close(db *gorm.DB) error {
	sdb, err := db.DB()
	if err != nil {
		return err
	}
	return sdb.Close()
}
