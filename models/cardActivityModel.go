package models

import (
	"gorm.io/gorm"
)

type CardActivity struct {
	gorm.Model
	CardNumber string `json:"cardNumber"`
	CardType   string `json:"cardType"`
	Event      string `json:"event"`
}
