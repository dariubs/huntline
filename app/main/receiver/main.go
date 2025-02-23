package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dariubs/go-producthunt"
	"github.com/dariubs/huntline/app/db"
	"github.com/dariubs/huntline/app/model"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var dbs *gorm.DB

func getYesterday() string {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	return yesterday.Format("2006-01-02")
}

func runAtScheduledTime(task func()) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

	for {
		now := time.Now().In(loc)
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 0, 30, 0, 0, loc)
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}
		durationUntilNextRun := time.Until(nextRun)
		log.Printf("Next run scheduled at: %s (PST)", nextRun)

		time.Sleep(durationUntilNextRun)
		task()
	}
}

func main() {
	// Define and parse command-line flag.
	runNow := flag.Bool("run-now", false, "Run the task immediately before starting the scheduler")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbs, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("PH_API_KEY")
	client := producthunt.ProductHunt{APIKey: apiKey}

	task := func() {
		date := getYesterday()
		products, err := client.GetProductsByRankByDate(date, 5)
		if err != nil {
			log.Fatalf("Error fetching products for date %s: %v", date, err)
		}

		fmt.Printf("Top Products on %s:\n", date)
		for i, product := range products {
			fmt.Printf("ID: %s\nName: %s\nTagline: %s\nWebsite: %s\nRank: %d\n\n",
				product.ID, product.Name, product.Tagline, product.Website, i)

			pdc := model.Product{
				Name:    product.Name,
				Tagline: product.Tagline,
				URL:     product.Website,
				Rank:    uint(i),
			}

			err = pdc.Save(dbs)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if *runNow {
		task()
	}

	runAtScheduledTime(task)
}
