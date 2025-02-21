package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dariubs/go-producthunt"
	"github.com/joho/godotenv"
)

// Product struct for storing API response
type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	Website     string `json:"website"`
}

func getYesterday() string {
	loc, _ := time.LoadLocation("America/Los_Angeles") // PST timezone
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	return yesterday.Format("2006-01-02")
}

func runAtScheduledTime(task func()) {
	loc, _ := time.LoadLocation("America/Los_Angeles")

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("PH_API_KEY")
	client := producthunt.ProductHunt{APIKey: apiKey}

	task := func() {
		date := getYesterday() // Get yesterday’s date
		products, err := client.GetProductsByRankByDate(date)
		if err != nil {
			log.Fatalf("Error fetching products for date %s: %v", date, err)
		}

		// Display the products
		fmt.Printf("Top Products on %s:\n", date)
		for _, product := range products {
			fmt.Printf("ID: %s\nName: %s\nTagline: %s\nWebsite: %s\n\n",
				product.ID, product.Name, product.Tagline, product.Website)
		}
	}

	// Run the scheduled task every day at 00:30 PST
	runAtScheduledTime(task)
}
