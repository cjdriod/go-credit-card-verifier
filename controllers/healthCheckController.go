package controllers

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timeStamp"`
}

func HealthCheckController(c *gin.Context) {
	data := HealthCheckResponse{
		Status:    configs.Constant.ApiStatus.Success,
		Message:   "Hello World!",
		TimeStamp: time.Now(),
	}

	c.JSON(http.StatusOK, data)
}
