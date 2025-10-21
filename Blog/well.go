package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONWell(c *gin.Context, status int, msg string) {
	rid := c.GetString("request_id")
	c.JSON(status, gin.H{
		"error":      gin.H{"message": msg},
		"request_id": rid,
	})
}

func WellCreateRequest(c *gin.Context, msg string) { JSONWell(c, http.StatusCreated, msg) }
func WellRequest(c *gin.Context, msg string)       { JSONWell(c, http.StatusOK, msg) }
