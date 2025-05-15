package services

import (
	"encoding/json"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/dto"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/repositories"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type CardServiceInterface interface {
	CreateCardFraudRecord(cardNumber string, cardType string, reason string) (*models.CardActivity, error)
	CheckCardValidity(cardNumber string) bool
	GetCardDetails(cardNumber string) *dto.CardInfoDto
	CreateUpdateBlackListRecord(userInfo models.User, cardNumber string, cardType string, banAction string) error
}

type CardService struct {
	cardRepo         repositories.CardRepositoryInterface
	blackListRepo    repositories.BlackListRepositoryInterface
	blackListService BlackListServiceInterface
}

func NewCardService(
	cardRepo repositories.CardRepositoryInterface,
	blackListRepo repositories.BlackListRepositoryInterface,
	blackListService BlackListServiceInterface,
) *CardService {
	return &CardService{
		cardRepo:         cardRepo,
		blackListRepo:    blackListRepo,
		blackListService: blackListService,
	}
}

func premiumCardDetailsLookup(cardNumber string) (*dto.BinCardResponse, error) {
	binUrl := fmt.Sprintf("https://lookup.binlist.net/%s", cardNumber[0:8])
	response, err := http.Get(binUrl)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Bin list api err:", response.StatusCode)
		return nil, err
	}

	var data dto.BinCardResponse
	err = json.NewDecoder(response.Body).Decode(&data)

	if err != nil {
		fmt.Println("Json parse err:", err.Error())
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	return &data, nil
}

func identifyCardNetworkBrand(cardNumber string) string {
	visaPattern := regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`)
	mastercardPattern := regexp.MustCompile(`^5[1-5][0-9]{14}$`)
	amexPattern := regexp.MustCompile(`^3[47][0-9]{13}$`)
	discoverPattern := regexp.MustCompile(`^6(?:011|5[0-9]{2})[0-9]{12}$`)
	jcbPattern := regexp.MustCompile(`^(?:2131|1800|35\d{3})\d{11}$`)
	dinersClubCard := regexp.MustCompile(`^3(?:0[0-5]|[68][0-9])[0-9]{11}$`)
	const (
		Visa       = "Visa"
		Mastercard = "Mastercard"
		Amex       = "American Express"
		Discover   = "Discover"
		Jcb        = "JCB"
		DinersClub = "Diners Club"
		Unknown    = "Unknown"
	)

	switch {
	case visaPattern.MatchString(cardNumber):
		return Visa
	case mastercardPattern.MatchString(cardNumber):
		return Mastercard
	case amexPattern.MatchString(cardNumber):
		return Amex
	case discoverPattern.MatchString(cardNumber):
		return Discover
	case jcbPattern.MatchString(cardNumber):
		return Jcb
	case dinersClubCard.MatchString(cardNumber):
		return DinersClub
	default:
		return Unknown
	}
}

func (service *CardService) CreateCardFraudRecord(cardNumber string, cardType string, reason string) (*models.CardActivity, error) {
	newCardActivity := &models.CardActivity{
		CardType:   cardType,
		CardNumber: cardNumber,
		Event:      reason,
	}

	if err := service.cardRepo.CreateCardActivity(newCardActivity); err != nil {
		return &models.CardActivity{}, fmt.Errorf("error creating fraud report %s", err)
	}

	return newCardActivity, nil
}

func (service *CardService) CheckCardValidity(cardNumber string) bool {
	if len(cardNumber) > 19 || len(cardNumber) < 14 {
		return false
	}

	// Transform string to individual int
	var cardNumberList []int
	for _, num := range cardNumber {
		num, err := strconv.Atoi(string(num))

		if err != nil {
			return false
		}
		cardNumberList = append(cardNumberList, num)
	}

	// Luhn algorithm validation
	sum, index, isSecond := 0, len(cardNumberList)-1, false
	for ; index >= 0; index-- {
		value := cardNumberList[index]

		if isSecond {
			value *= 2

			if value > 9 {
				value -= 9
			}
		}

		sum += value
		isSecond = !isSecond
	}
	return sum%10 == 0
}

func (service *CardService) GetCardDetails(cardNumber string) *dto.CardInfoDto {
	cardInfoDto := &dto.CardInfoDto{
		IsBlackListed: service.blackListService.HasBlackListRecord(cardNumber),
		Network:       identifyCardNetworkBrand(cardNumber),
	}

	if !cardInfoDto.IsBlackListed && configs.Constant.EnablePremiumCardCheck {
		if response, err := premiumCardDetailsLookup(cardNumber); err == nil {
			if response.Type != "" {
				cardInfoDto.Type = response.Type
			}

			if response.Brand != "" {
				cardInfoDto.Network = response.Brand
			}

			if response.Bank.Name != "" {
				cardInfoDto.Issuer = response.Bank.Name
			}

			if response.Country.Name != "" {
				cardInfoDto.IssueCountry = response.Country.Name
			}
		}

	}

	if activities, err := service.cardRepo.GetCardActivities(cardNumber); err == nil {
		cardInfoDto.Activities = activities
	}

	return cardInfoDto
}

func (service *CardService) CreateUpdateBlackListRecord(adminInfo models.User, cardNumber string, cardType string, banAction string) error {
	blackListRecord, blackListDbError := service.blackListRepo.FindBlackListRecordByCardNumber(cardNumber)
	if blackListDbError != nil {
		return fmt.Errorf("error getting black list record")
	}

	transaction := database.DB.Begin()
	if err := transaction.Error; err != nil {
		return fmt.Errorf("fail to start transaction: %s", err)
	}

	var hasChanges bool
	var eventMessage string
	if banAction == configs.Constant.CardBlackListAction.Ban && blackListRecord.CardNumber == "" {
		newBlackListRecord := &models.BlackList{
			CardNumber:  cardNumber,
			DeletedById: adminInfo.ID,
		}

		if err := service.blackListRepo.CreateBlackListRecord(newBlackListRecord, transaction); err != nil {
			transaction.Rollback()
			return fmt.Errorf("error creating black list record: %s", err)
		}

		eventMessage = fmt.Sprintf("Event: Card BAN since %s", time.Now().UTC())
		hasChanges = true
	}

	if banAction == configs.Constant.CardBlackListAction.Unban && blackListRecord.CardNumber == cardNumber {
		if err := service.blackListRepo.DeleteBlackListRecord(cardNumber, transaction); err != nil {
			transaction.Rollback()
			return fmt.Errorf("error deleting black list record: %s", err)
		}

		eventMessage = fmt.Sprintf("Event: Card un-BAN since %s", time.Now().UTC())
		hasChanges = true
	}

	if hasChanges {
		if eventMessage == "" {
			transaction.Rollback()
			return fmt.Errorf("error creating card activity: event message is empty")
		}

		cardActivityPayload := &models.CardActivity{
			CardType:   cardType,
			CardNumber: cardNumber,
			Event:      eventMessage,
		}

		if err := service.cardRepo.CreateCardActivity(cardActivityPayload); err != nil {
			transaction.Rollback()
			return fmt.Errorf("error creating card activity: %s", err)
		}

		if err := transaction.Commit().Error; err != nil {
			return fmt.Errorf("error committing black list record: %s", err)
		}
	}

	return nil
}
