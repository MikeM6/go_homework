log Service

A simple blog API with user registration/login (JWT), posts (CRUD), and comments.
Built with gin + gorm + MySQL.
Prerequisites

Go 1.20+ (recommended)
MySQL 5.7+/8.0+
Git, curl (for testing)
Optional: Postman/Insomnia
Tech Stack

Web: github.com/gin-gonic/gin
ORM: gorm.io/gorm + gorm.io/driver/mysql
Auth: github.com/golang-jwt/jwt/v5
Password: golang.org/x/crypto/bcrypt
Project Layout

Code: Blog/
main.go: server bootstrap and route registration
dbModel.go: models and DB init
loginAndRegistration.go: auth (register/login/JWT/middleware)
postManagement.go: posts CRUD
commentManagement.go: comments
errors.go (if present): unified error helpers
middleware.go (if present): request ID, recover, access logs
Database Setup

Create a database (default name is gorm per DSN):
CREATE DATABASE gorm CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
Ensure a user with access:
Default DSN in code is root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local
If your credentials/db differ, update the DSN in Blog/main.go or adapt to an environment variable.
Environment Variables

BLOG_JWT_SECRET (recommended): secret key for signing JWT. Example:
Linux/macOS: export BLOG_JWT_SECRET="change_me"
Windows PowerShell: $env:BLOG_JWT_SECRET="change_me"
Install & Run

In a new terminal:
cd Blog
go mod tidy
Start MySQL locally and ensure the DSN is valid.
Run the server:
go run .
The API listens on http://localhost:8080
Auto-migration:
Tables for users, posts, comments are auto-created on startup.
API Endpoints

Auth
POST /api/register — register user: { "username": "...", "email": "...", "password": "..." }
POST /api/login — login via username or email: { "username": "...", "email": "...", "password": "..." }
GET /api/me — requires Authorization: Bearer <token>
Posts
POST /api/posts — create (auth)
GET /api/posts — list
GET /api/posts/:id — detail
PUT /api/posts/:id — update (author only)
DELETE /api/posts/:id — delete (author only)
Comments
POST /api/posts/comments — create (auth): { "postid": 1, "content": "..." }
GET /api/posts/:id/comments — list for a post
Quick Test (curl)

Register:
curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d "{\"username\":\"u1\",\"email\":\"u1@example.com\",\"password\":\"secret123\"}"
Login (get token):
curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\"username\":\"u1\",\"password\":\"secret123\"}"
Create post:
curl -X POST http://localhost:8080/api/posts -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d "{\"title\":\"t\",\"content\":\"c\"}"
Logging & Error Handling

Handlers return consistent HTTP codes: 400/401/403/404/500 with JSON error bodies.
If errors.go/middleware.go are in place:
Unified JSON errors with request_id.
Access logs per request, DB errors surfaced via c.Error(err) and included in logs.
Panic recovery returns 500 and logs stack traces.
GORM logs at Info level to stdout.
Troubleshooting

Cannot connect to DB:
Verify MySQL is running, credentials match DSN, and DB gorm exists.
401 Unauthorized:
Ensure Authorization: Bearer <token> is set; token not expired.
403 Forbidden:
Only the post author can update/delete.
Port in use:
Change the port in Blog/main.go (r.Run(":8080")).