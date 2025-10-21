package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 3. 用户认证与授权
// 实现用户注册和登录功能，用户注册时需要对密码进行加密存储，登录时验证用户输入的用户名和密码。
// 使用 JWT（JSON Web Token）实现用户认证和授权，用户登录成功后返回一个 JWT，后续的需要认证的接口需要验证该 JWT 的有效性。

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// LoginRequest 登录请求体（支持用户名或邮箱）
type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// AuthResponse 通用响应（登录/注册）
type AuthResponse struct {
	Token    string `json:"token,omitempty"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// JWTClaims 自定义 JWT Claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// getJWTSecret 从环境变量读取密钥，或使用默认值（仅用于开发）
func getJWTSecret() []byte {
	if v := os.Getenv("BLOG_JWT_SECRET"); v != "" {
		return []byte(v)
	}
	return []byte("dev_secret_change_me")
}

// hashPassword 使用 bcrypt 对密码进行哈希
func hashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// verifyPassword 校验密码
func verifyPassword(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// generateToken 生成 JWT token
func generateToken(user *User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user_auth",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// parseToken 验证并解析 token
func parseToken(tokenStr string) (*JWTClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*JWTClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// RegisterAuthRoutes 注册认证相关路由
//
// r: Gin engine
// db: GORM 实例
// base: 路由前缀，例如 "/api"
func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB, base string) {
	g := r.Group(base)

	// 用户注册
	g.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			BadRequest(c, err.Error())
			return
		}

		// 唯一性检查
		var cnt int64
		if err := db.Model(&User{}).Where("username = ?", req.Username).Or("email = ?", req.Email).Count(&cnt).Error; err != nil {
			c.Error(err)
			Internal(c, "database error : get user count fail")
			return
		}
		if cnt > 0 {
			WellRequest(c, "username or email already exists")
			return
		}

		hashed, err := hashPassword(req.Password)
		if err != nil {
			c.Error(err)
			Internal(c, "failed to hash password")
			return
		}

		user := &User{Username: req.Username, Email: req.Email, Password: hashed}
		if err := db.Create(user).Error; err != nil {
			c.Error(err)
			Internal(c, "failed to create user")
			return
		}

		c.JSON(http.StatusCreated, AuthResponse{UserID: user.ID, Username: user.Username, Email: user.Email})
	})

	// 用户登录
	g.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			BadRequest(c, err.Error())
			return
		}

		if req.Username == "" && req.Email == "" {
			BadRequest(c, "username or email is required")
			return
		}

		var user User
		q := db
		if req.Username != "" {
			q = q.Where("username = ?", req.Username)
		}
		if req.Email != "" {
			q = q.Or("email = ?", req.Email)
		}
		if err := q.First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				WellRequest(c, "invalid credentials")
				return
			}
			c.Error(err)
			Internal(c, "database error")
			return
		}

		if err := verifyPassword(user.Password, req.Password); err != nil {
			WellRequest(c, "invalid credentials")
			return
		}

		token, err := generateToken(&user, 24*time.Hour)
		if err != nil {
			c.Error(err)
			Internal(c, "failed to create token")
			return
		}
		c.JSON(http.StatusOK, AuthResponse{Token: token, UserID: user.ID, Username: user.Username, Email: user.Email})
	})

	// 示例受保护接口：获取当前用户信息
	g.GET("/me", AuthMiddleware(), func(c *gin.Context) {
		uid, _ := c.Get("currentUserID")
		uname, _ := c.Get("currentUsername")
		c.JSON(http.StatusOK, gin.H{"user_id": uid, "username": uname})
	})
}
