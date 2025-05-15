package router

import (
	"github.com/cjdriod/go-credit-card-verifier/controllers"
	"github.com/cjdriod/go-credit-card-verifier/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupCardRoutes(router *gin.RouterGroup, controllers *controllers.CardController) {
	v1 := router.Group("/v1/card")

	v1.GET("/:cardNumber", controllers.GetCardInfo)
	v1.POST("/report-fraud", controllers.ReportCardFraud)
	v1.POST("/black-list", middlewares.AuthenticationMiddleware(), controllers.BlackListCard)
}
