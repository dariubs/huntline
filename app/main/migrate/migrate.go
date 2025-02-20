package main

import (
	"log"

	"github.com/dariubs/huntline/app/model"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = model.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}
}
