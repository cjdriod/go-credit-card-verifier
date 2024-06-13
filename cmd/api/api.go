package api

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/controllers"
	"github.com/cjdriod/go-credit-card-verifier/middlewares"
	"github.com/gin-gonic/gin"
	"log"
)

type RestServer struct {
	router *gin.Engine
}
type Card struct {
	Network    string `json:"network" binding:"required"`
	CardNumber string `json:"card_number" binding:"required"`
}

func InitServer() *RestServer {
	r := gin.Default()
	r.GET("/hello-world", controllers.HealthCheckController)

	// Card Route
	r.GET("/card/:cardNumber", controllers.GetCardInfoController)
	r.POST("/card/report-fraud", controllers.ReportCardFraudController)
	r.POST("/card/block", middlewares.RequireAuthentication, controllers.BlackListCardController)

	// User Route
	r.POST("/user/register", controllers.RegisterUserController)
	r.POST("/user/login", controllers.LoginUserController)
	r.POST("/user/logout", middlewares.RequireAuthentication, controllers.LogoutUserController)

	return &RestServer{router: r}
}

func (s RestServer) Serve() {
	if err := s.router.Run(configs.Constant.Host); err != nil {
		log.Fatal(err)
	}
}
