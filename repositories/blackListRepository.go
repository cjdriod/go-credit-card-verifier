package repositories

import (
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"gorm.io/gorm"
)

type BlackListRepositoryInterface interface {
	CreateBlackListRecord(record *models.BlackList, transaction *gorm.DB) error
	FindBlackListRecordByCardNumber(cardNumber string) (*models.BlackList, error)
	DeleteBlackListRecord(cardNumber string, transaction *gorm.DB) error
}
type BlackListRepository struct{}

func NewBlackListRepository() *BlackListRepository {
	return &BlackListRepository{}
}

func (repo *BlackListRepository) CreateBlackListRecord(record *models.BlackList, transaction *gorm.DB) error {
	dbInstance := database.GetDB(transaction)
	return dbInstance.Create(record).Error
}

func (repo *BlackListRepository) FindBlackListRecordByCardNumber(cardNumber string) (*models.BlackList, error) {
	result := &models.BlackList{}
	err := database.DB.Where("card_number = ?", cardNumber).First(&result).Error

	return result, err
}

func (repo *BlackListRepository) DeleteBlackListRecord(cardNumber string, transaction *gorm.DB) error {
	dbInstance := database.GetDB(transaction)
	return dbInstance.Delete(&models.BlackList{CardNumber: cardNumber}).Error
}
