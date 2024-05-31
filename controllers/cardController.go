package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/Utils"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/database"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Card struct {
	Type          string   `json:"type"`
	Issuer        string   `json:"issuer"`
	IssueCountry  string   `json:"issueCountry"`
	Network       string   `json:"network"`
	CardNumber    string   `json:"cardNumber"`
	Activities    []string `json:"activities"`
	IsBlackListed bool     `json:"isBlackListed"`
}

type BinCardResponse struct {
	Brand   string `json:"brand"`
	Type    string `json:"type"`
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
	Bank struct {
		Name string `json:"name"`
	} `json:"bank"`
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

func validateCardNumber(cardNumber string) bool {
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

func premiumCardDetailsLookup(cardNumber string) (*BinCardResponse, error) {
	binUrl := fmt.Sprintf("https://lookup.binlist.net/%s", cardNumber[7:8])
	response, err := http.Get(binUrl)

	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Bin list api err:", response.StatusCode)
		return nil, err
	}

	var data BinCardResponse
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

func GetCardInfoController(c *gin.Context) {
	invalidBankCardErr := errors.New("invalid bank card number")
	cardNumber := c.Param("cardNumber")

	if len(cardNumber) > 19 || len(cardNumber) < 14 {
		Utils.BadRequest(c, []string{invalidBankCardErr.Error()})
		return
	}

	var (
		isBlackListCard  = false
		isValidCard      = validateCardNumber(cardNumber)
		activities       = make([]string, 0)
		cardType         = ""
		cardIssueCountry = ""
		cardIssuer       = ""
		cardNetwork      = identifyCardNetworkBrand(cardNumber)
	)

	if !isValidCard {
		Utils.BadRequest(c, []string{invalidBankCardErr.Error()})
		return
	}

	blackListQry := database.DB.First(&models.BlackList{CardNumber: cardNumber})
	isBlackListCard = blackListQry.RowsAffected != 0

	if !isBlackListCard && configs.Constant.EnablePremiumCardCheck {
		response, err := premiumCardDetailsLookup(cardNumber)
		if err == nil {
			if response.Type != "" {
				cardType = response.Type
			}

			if response.Brand != "" {
				cardNetwork = response.Brand
			}

			if response.Bank.Name != "" {
				cardIssuer = response.Bank.Name
			}

			if response.Country.Name != "" {
				cardIssueCountry = response.Country.Name
			}
		}
	}

	var activityRecords []models.CardActivity
	database.DB.Select("event").Order("updated_at desc").Find(&activityRecords)
	for _, activity := range activityRecords {
		activities = append(activities, activity.Event)
	}

	message := "Valid (debit/credit) card number"
	switch {
	case isBlackListCard:
		message = "This is a black listed credit/debit card"

	case cardNetwork == "Unknown":
		message = "Unable to detect card brand at the moment"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": Card{
			Type:          cardType,
			Issuer:        cardIssuer,
			Network:       cardNetwork,
			CardNumber:    cardNumber,
			Activities:    activities,
			IssueCountry:  cardIssueCountry,
			IsBlackListed: isBlackListCard,
		},
		"status":  configs.Constant.ApiStatus.Success,
		"message": message,
	})
}

func ReportCardFraudController(c *gin.Context) {
	reportFailedErr := errors.New("failed to file fraud report")
	var body struct {
		CardNumber string `json:"cardNumber"`
		CardType   string `json:"cardType"`
		Reason     string `json:"reason"`
	}

	err := c.BindJSON(&body)
	if err != nil || body.CardNumber == "" || body.Reason == "" {
		Utils.BadRequest(c, []string{reportFailedErr.Error()})
		return
	}
	fmt.Println(body)
	fraudReport := models.CardActivity{CardType: body.CardType, CardNumber: body.CardNumber, Event: body.Reason}
	result := database.DB.Create(&fraudReport)

	if result.Error != nil {
		fmt.Println("Database error:", result.Error)
		Utils.BadRequest(c, []string{reportFailedErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    body,
		"status":  configs.Constant.ApiStatus.Success,
		"message": "Thank you for reporting abnormal card activity with our system, we will review soon",
	})

}

func BlackListCardController(c *gin.Context) {
	reportFailedErr := errors.New("failed to black list card")
	var body struct {
		CardNumber string `json:"cardNumber"`
		CardType   string `json:"cardType"`
		IsChecked  bool   `json:"IsChecked"`
	}

	user, err := c.Get("user")
	if bindErr := c.BindJSON(&body); !err || bindErr != nil || body.CardNumber == "" {
		Utils.BadRequest(c, []string{reportFailedErr.Error()})
		return
	}

	var (
		blackList  models.BlackList
		hasChanges = false
	)
	databaseErr := database.DB.Transaction(func(tx *gorm.DB) error {
		database.DB.First(&blackList)

		switch {
		case body.IsChecked && blackList.CardNumber == "":
			if err := tx.Create(&models.BlackList{
				CardNumber:  body.CardNumber,
				DeletedById: user.(models.User).ID,
			}).Error; err != nil {
				return err
			}
			hasChanges = true

		//case body.IsChecked && blackList.CardNumber == body.CardNumber :
		//	if err := tx.Model(&models.BlackList{}).Where("card_number = ?", body.CardNumber).Updates(models.BlackList{
		//		DeletedById: user.(models.User).ID,
		//	}).Error; err != nil {
		//		return err
		//	}
		//	hasChanges = true

		case !body.IsChecked && blackList.CardNumber == body.CardNumber:
			if err := tx.Delete(&models.BlackList{CardNumber: body.CardNumber}).Error; err != nil {
				return err
			}
			hasChanges = true
		}

		if hasChanges {
			var event string
			if body.IsChecked {
				event = fmt.Sprintf("Event: Card BAN since %v", time.Now().UTC())
			} else {
				event = fmt.Sprintf("Event: Card un-BAN since %v", time.Now().UTC())
			}

			if err := tx.Create(&models.CardActivity{CardType: body.CardType, CardNumber: body.CardNumber, Event: event}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if databaseErr != nil {
		fmt.Println("Database error:", databaseErr)
		Utils.BadRequest(c, []string{reportFailedErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("This card (%v)is recorded successfully", body.CardNumber),
	})
}
