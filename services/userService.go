package services

import (
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserServiceInterface interface {
	CreateNewUser(username string, password string, email string) error
	LoginUser(email string, password string) (string, string, error)
}

type UserService struct {
	UserRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (userService *UserService) CreateNewUser(username string, password string, email string) error {
	hashPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), 10)

	if hashErr != nil {
		return fmt.Errorf("error hashing password: %s", hashErr)
	}

	newUser := &models.User{
		Email:    email,
		Username: username,
		Password: string(hashPassword),
	}

	if err := userService.UserRepo.CreateUser(newUser); err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	return nil
}

func (userService *UserService) LoginUser(email string, password string) (string, string, error) {
	user, userErr := userService.UserRepo.FindUserByEmail(email)

	if userErr != nil || user.ID == 0 {
		return "", "", fmt.Errorf("error finding user by email: %s", email)
	}

	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); bcryptErr != nil {
		return "", "", fmt.Errorf("error hashing password: %s", bcryptErr)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, signErr := token.SignedString(configs.Constant.JwtSecret)

	if signErr != nil {
		return "", "", fmt.Errorf("error signing token: %s", signErr)
	}

	return user.Username, tokenString, nil

}
