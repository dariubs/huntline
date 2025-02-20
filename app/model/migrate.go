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
	// Auto migrate  models
	err = DB.AutoMigrate(&Product{})
	if err != nil {
		return err
	}

	return nil
}
