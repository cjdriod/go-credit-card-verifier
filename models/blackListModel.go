package models

import (
	"time"
)

type BlackList struct {
	CardNumber  string `gorm:"primaryKey"`
	DeletedById uint   `gorm:"index"`
	DeletedBy   User   `gorm:"foreignKey:DeletedById"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
