package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	gorm.Model

	Name        string `gorm:"type:varchar(255);not null"`
	URL         string `gorm:"type:text;not null"`
	Tagline     string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Rank        uint
	Logo        string    `gorm:"type:text"`
	Date        time.Time `gorm:"type:date"`
}

func (product *Product) Save(db *gorm.DB) error {
	product.Date = product.Date.Truncate(24 * time.Hour)

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "date"}}, // Conflict based on name and date
		DoUpdates: clause.Assignments(map[string]interface{}{"rank": product.Rank}),
	}).Create(product).Error

	if err != nil {
		return err
	}

	return nil
}
