package database

import (
	"github.com/cjdriod/go-credit-card-verifier/configs"
	"github.com/cjdriod/go-credit-card-verifier/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(mysql.Open(configs.Constant.MySqlConnectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect mySql database")
	}
}

func SyncDatabase() {
	var err error

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to initialized users table")
	}

	err = DB.AutoMigrate(&models.BlackList{})
	if err != nil {
		panic("failed to initialized blackLists table")
	}

	err = DB.AutoMigrate(&models.CardActivity{})
	if err != nil {
		panic("failed to initialized cardActivity table")
	}
}
