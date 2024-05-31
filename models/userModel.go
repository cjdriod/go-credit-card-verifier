package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);unique;index;not null"`
	Username string `json:"username"`
	Password string `json:"-"`
}
