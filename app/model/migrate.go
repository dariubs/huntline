package model

import (
	"github.com/dariubs/huntline/app/db"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AutoMigrate() error {
	DB, err := db.ConnectToDB()
	if err != nil {
		return err
	}
	// Auto migrate models
	err = DB.AutoMigrate(&Product{})
	if err != nil {
		return err
	}

	// Migrate existing data: set platform="producthunt" for products without a platform
	result := DB.Model(&Product{}).Where("platform = '' OR platform IS NULL").Update("platform", "producthunt")
	if result.Error != nil {
		return result.Error
	}

	return nil
}
