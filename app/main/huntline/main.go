package main

import (
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"

	"github.com/dariubs/huntline/app/db"
	"github.com/dariubs/huntline/app/handler/huntline"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// extractDomain extracts the domain from a URL string
func extractDomain(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// If parsing fails, try to extract domain manually
		if len(urlStr) > 7 && (urlStr[:7] == "http://" || urlStr[:8] == "https://") {
			start := 0
			if urlStr[:7] == "http://" {
				start = 7
			} else {
				start = 8
			}
			end := len(urlStr)
			if idx := findChar(urlStr, start, '/'); idx > 0 {
				end = idx
			}
			if idx := findChar(urlStr, start, '?'); idx > 0 && idx < end {
				end = idx
			}
			if idx := findChar(urlStr, start, '#'); idx > 0 && idx < end {
				end = idx
			}
			return urlStr[start:end]
		}
		return urlStr
	}
	return parsedURL.Hostname()
}

func findChar(s string, start int, c byte) int {
	for i := start; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

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

	// Register template functions
	router.SetFuncMap(template.FuncMap{
		"extractDomain": extractDomain,
	})

	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("view/huntline/**/*")

	router.GET("/", huntline.IndexHandler(dbs, gd))
	router.GET("/api/timeline", huntline.IndexAPIHandler(dbs, gd))
	router.GET("/archive", huntline.ArchiveHandler(dbs, gd))
	router.GET("/best/month", huntline.BestMonthHandler(dbs, gd))
	router.GET("/best/week", huntline.BestWeekHandler(dbs, gd))
	router.GET("/platforms", huntline.PlatformsHandler(dbs, gd))

	port := os.Getenv("HL_PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf("0.0.0.0:%s", port))

}
