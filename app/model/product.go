package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string `gorm:"type:varchar(255);not null"`
	URL         string `gorm:"type:text";not null`
	Tagline     string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Rank        uint
	Logo        string    `gorm:"type:text"`
	Date        time.Time `gorm:"type:date"`
}
