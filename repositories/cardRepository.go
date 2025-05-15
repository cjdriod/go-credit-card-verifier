package repositories

import (
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
)

type CardRepositoryInterface interface {
	CreateCardActivity(activity *models.CardActivity) error
	GetCardActivities(cardNumber string) ([]string, error)
}
type CardRepository struct{}

func NewCardRepository() *CardRepository {
	return &CardRepository{}
}

func (repo *CardRepository) CreateCardActivity(activity *models.CardActivity) error {
	return database.DB.Create(activity).Error
}

func (repo *CardRepository) GetCardActivities(cardNumber string) ([]string, error) {
	var cardActivities []string

	err := database.DB.
		Model(&models.CardActivity{}).
		Where("CardNumber = ?", cardNumber).
		Pluck("Event", &cardActivities).
		Error

	return cardActivities, err
}
