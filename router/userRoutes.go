package router

import (
	"github.com/cjdriod/go-credit-card-verifier/controllers"
	"github.com/cjdriod/go-credit-card-verifier/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(routerGroup *gin.RouterGroup, userControllers *controllers.UserController) {
	v1 := routerGroup.Group("/v1/user")

	v1.POST("/register", middlewares.AuthenticationMiddleware(), userControllers.RegisterNewUser)
	v1.POST("/login", userControllers.LoginUser)
	v1.POST("/logout", middlewares.AuthenticationMiddleware(), userControllers.LogoutUser)
}
