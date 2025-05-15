package repositories

import (
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func (repo *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.Where("email = ?", email).First(user).Error

	return user, err
}
