package main

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware Gin 中间件：校验 Authorization: Bearer <token>
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if hdr == "" || !strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
			Unauthorized(c, "missing or invalid authorization header")
			return
		}
		tokenStr := strings.TrimSpace(hdr[len("Bearer "):])
		claims, err := parseToken(tokenStr)
		if err != nil {
			c.Error(err)
			Unauthorized(c, "invalid token")
			return
		}
		// 将用户信息放入上下文，便于后续处理
		c.Set("currentUserID", claims.UserID)
		c.Set("currentUsername", claims.Username)
		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = time.Now().UTC().Format("20060102T150405.000000000Z07:00")
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic rid=%s err=%v stack=%s", c.GetString("request_id"), r, debug.Stack())
				Internal(c, "internal server error")
			}
		}()
		c.Next()
	}
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		uid, _ := c.Get("currentUserID")
		log.Printf(
			"rid=%s status=%d method=%s path=%s ip=%s dur=%s uid=%v err=%s",
			c.GetString("request_id"),
			c.Writer.Status(),
			c.Request.Method,
			c.FullPath(),
			c.ClientIP(),
			dur,
			uid,
			c.Errors.ByType(gin.ErrorTypeAny).String(),
		)
	}
}
