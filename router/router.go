package router

import (
	"github.com/cjdriod/go-credit-card-verifier/controllers"
	"github.com/cjdriod/go-credit-card-verifier/dependencyInjection"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	diControllers := dependencyInjection.InitDependencies()

	api.GET("/hello-world", controllers.HealthCheckController)

	SetupUserRoutes(api, diControllers.UserController)
	SetupCardRoutes(api, diControllers.CardController)

	return router
}
