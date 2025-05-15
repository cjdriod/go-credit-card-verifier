package controllers

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HealthCheckController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    configs.Constant.ApiStatus.Success,
		"message":   "Hello World!",
		"timeStamp": time.Now(),
	})
}
