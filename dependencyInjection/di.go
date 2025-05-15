package dependencyInjection

import (
	"github.com/cjdriod/go-credit-card-verifier/controllers"
	"github.com/cjdriod/go-credit-card-verifier/repositories"
	"github.com/cjdriod/go-credit-card-verifier/services"
)

type Di struct {
	UserController *controllers.UserController
	CardController *controllers.CardController
}

func InitDependencies() *Di {

	// Repos
	userRepo := repositories.NewUserRepository()
	cardRepo := repositories.NewCardRepository()
	blackListRepo := repositories.NewBlackListRepository()

	// services
	userService := services.NewUserService(userRepo)
	blackListService := services.NewBlackListService(blackListRepo)
	cardService := services.NewCardService(cardRepo, blackListRepo, blackListService)

	// controllers
	userController := controllers.NewUserController(userService)
	cardController := controllers.NewCardController(cardService)

	return &Di{
		UserController: userController,
		CardController: cardController,
	}
}
