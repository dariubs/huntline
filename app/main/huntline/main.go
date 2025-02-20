package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dariubs/huntline/app/db"
	"github.com/dariubs/huntline/app/handler/huntline"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var dbs *gorm.DB
var gd types.General

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbs, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	gd = types.General{
		Name:    os.Getenv("HL_NAME"),
		Logo:    os.Getenv("HL_LOGO"),
		Favicon: os.Getenv("HL_FAVICON"),
		URL:     os.Getenv("HL_URL"),
		CDN:     os.Getenv("HL_CDN"),
		X:       os.Getenv("HL_X"),
		GitHub:  os.Getenv("HL_GITHUB"),
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Delims("{{", "}}")

	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("view/huntline/**/*")

	router.GET("/", huntline.IndexHandler(dbs, gd))

	router.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("HL_PORT")))

}
