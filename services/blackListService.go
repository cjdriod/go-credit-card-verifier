package services

import (
	"github.com/cjdriod/go-credit-card-verifier/repositories"
)

type BlackListServiceInterface interface {
	HasBlackListRecord(cardNumber string) bool
}
type BlackListService struct {
	blackListRepo repositories.BlackListRepositoryInterface
}

func NewBlackListService(blackListRepo repositories.BlackListRepositoryInterface) *BlackListService {
	return &BlackListService{blackListRepo: blackListRepo}
}

func (service *BlackListService) HasBlackListRecord(cardNumber string) bool {
	record, err := service.blackListRepo.FindBlackListRecordByCardNumber(cardNumber)

	return err == nil && record != nil
}
