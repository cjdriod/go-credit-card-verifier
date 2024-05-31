package Utils

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadRequest(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": errors,
		"status": configs.Constant.ApiStatus.Fail,
	})
}

func UnauthorizedRequest(c *gin.Context, errors []string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"errors": errors,
		"status": configs.Constant.ApiStatus.Fail,
	})
}
