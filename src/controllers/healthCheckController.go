package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": os.Getenv("APP_NAME"),
		"version": os.Getenv("APP_VERSION"),
	})
}
