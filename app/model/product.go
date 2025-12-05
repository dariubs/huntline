package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_date_platform"`
	URL         string `gorm:"type:text;not null"`
	Tagline     string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Rank        uint
	Logo        string    `gorm:"type:text"`
	Date        time.Time `gorm:"type:date;uniqueIndex:idx_name_date_platform"`
	Platform    string    `gorm:"type:varchar(100);not null;default:'producthunt';uniqueIndex:idx_name_date_platform"`
}

func (product *Product) Save(db *gorm.DB) error {
	product.Date = product.Date.Truncate(24 * time.Hour)
	
	// Set default platform if not specified
	if product.Platform == "" {
		product.Platform = "producthunt"
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "date"}, {Name: "platform"}}, // Conflict based on name, date, and platform
		DoUpdates: clause.Assignments(map[string]interface{}{"rank": product.Rank}),
	}).Create(product).Error

	if err != nil {
		return err
	}

	return nil
}
