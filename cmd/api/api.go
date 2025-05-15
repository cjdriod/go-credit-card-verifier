package api

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/router"
	"github.com/gin-gonic/gin"
	"log"
)

type RestServer struct {
	router *gin.Engine
}

func InitServer() *RestServer {
	r := router.SetupRouter()

	return &RestServer{router: r}
}

func (s RestServer) Serve() {
	if err := s.router.Run(configs.Constant.Host); err != nil {
		log.Fatal(err)
	}
}
