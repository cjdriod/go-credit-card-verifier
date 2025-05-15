package controllers

import (
	"errors"
	"fmt"
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"github.com/cjdriod/go-credit-card-verifier/services"
	"github.com/cjdriod/go-credit-card-verifier/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CardController struct {
	cardService services.CardServiceInterface
}

func NewCardController(cardService services.CardServiceInterface) *CardController {
	return &CardController{cardService: cardService}
}

func (controller *CardController) ReportCardFraud(c *gin.Context) {
	reportFailedErr := errors.New("failed to file fraud report")
	var body struct {
		CardNumber string `json:"cardNumber"`
		CardType   string `json:"cardType"`
		Reason     string `json:"reason"`
	}

	if err := c.BindJSON(&body); err != nil || body.CardNumber == "" || body.Reason == "" || body.CardType == "" {
		utils.BadRequest(c, []string{reportFailedErr.Error()})
	}

	newFraudRecord, err := controller.cardService.CreateCardFraudRecord(body.CardNumber, body.CardType, body.Reason)

	if err != nil {
		utils.BadRequest(c, []string{reportFailedErr.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    newFraudRecord,
		"status":  configs.Constant.ApiStatus.Success,
		"message": "Thank you for reporting abnormal card activity with our system, we will review soon",
	})
}

func (controller *CardController) GetCardInfo(c *gin.Context) {
	invalidBankCardErr := errors.New("invalid bank card number")
	cardNumber := c.Param("cardNumber")

	if !controller.cardService.CheckCardValidity(cardNumber) {
		utils.BadRequest(c, []string{invalidBankCardErr.Error()})
		return
	}

	cardDetails := controller.cardService.GetCardDetails(cardNumber)

	var message string
	if cardDetails.IsBlackListed {
		message = "This is a black listed credit/debit card"
	} else if cardDetails.Network == "" {
		message = "Unable to detect card brand at the moment"
	} else {
		message = fmt.Sprintf("%v card activitie(s) found", len(cardDetails.Activities))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    cardDetails,
		"status":  configs.Constant.ApiStatus.Success,
		"message": message,
	})
}

func (controller *CardController) BlackListCard(c *gin.Context) {
	reportFailedErr := errors.New("failed to black list card")
	var body struct {
		CardNumber  string `json:"cardNumber"`
		CardType    string `json:"cardType"`
		IsBanAction bool   `json:"IsBanAction"`
	}

	userInfo, userErr := c.Get("user")
	if bindErr := c.BindJSON(&body); !userErr || bindErr != nil || body.CardNumber == "" {
		utils.BadRequest(c, []string{reportFailedErr.Error()})
		return
	}

	banAction := configs.Constant.CardBlackListAction.Unban
	if body.IsBanAction {
		banAction = configs.Constant.CardBlackListAction.Ban
	}

	if err := controller.cardService.CreateUpdateBlackListRecord(
		userInfo.(models.User),
		body.CardNumber,
		body.CardType,
		banAction,
	); err != nil {
		utils.BadRequest(c, []string{reportFailedErr.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  configs.Constant.ApiStatus.Success,
		"message": fmt.Sprintf("This card (%s)is recorded successfully", body.CardNumber),
	})
}
