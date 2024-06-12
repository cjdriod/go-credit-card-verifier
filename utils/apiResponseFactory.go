package utils

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func genericErrorHandler(c *gin.Context, errors []string, errorCode int) {
	c.JSON(errorCode, gin.H{
		"errors": errors,
		"status": configs.Constant.ApiStatus.Fail,
	})
}

func BadRequest(c *gin.Context, errors []string) {
	genericErrorHandler(c, errors, http.StatusBadRequest)
}

func UnauthorizedRequest(c *gin.Context, errors []string) {
	genericErrorHandler(c, errors, http.StatusUnauthorized)
}
