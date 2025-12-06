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
	// Normalize date to midnight in San Francisco timezone to avoid timezone conversion issues
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err == nil {
		// Convert to PST/PDT timezone first, then normalize to midnight
		dateInLoc := product.Date.In(loc)
		product.Date = time.Date(dateInLoc.Year(), dateInLoc.Month(), dateInLoc.Day(), 0, 0, 0, 0, loc)
	} else {
		// Fallback: truncate to 24 hours if timezone loading fails
		product.Date = product.Date.Truncate(24 * time.Hour)
	}
	
	// Set default platform if not specified
	if product.Platform == "" {
		product.Platform = "producthunt"
	}

	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}, {Name: "date"}, {Name: "platform"}}, // Conflict based on name, date, and platform
		DoUpdates: clause.Assignments(map[string]interface{}{
			"rank":        product.Rank,
			"tagline":     product.Tagline,
			"url":         product.URL,
			"logo":        product.Logo,
			"description": product.Description,
		}),
	}).Create(product).Error

	if err != nil {
		return err
	}

	return nil
}
