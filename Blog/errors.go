package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, msg string) {
	rid := c.GetString("request_id")
	c.AbortWithStatusJSON(status, gin.H{
		"error":      gin.H{"message": msg},
		"request_id": rid,
	})
}

func BadRequest(c *gin.Context, msg string)   { JSONError(c, http.StatusBadRequest, msg) }
func Unauthorized(c *gin.Context, msg string) { JSONError(c, http.StatusUnauthorized, msg) }
func Forbidden(c *gin.Context, msg string)    { JSONError(c, http.StatusForbidden, msg) }
func NotFound(c *gin.Context, msg string)     { JSONError(c, http.StatusNotFound, msg) }
func Internal(c *gin.Context, msg string)     { JSONError(c, http.StatusInternalServerError, msg) }
