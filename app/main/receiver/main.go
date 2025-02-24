package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dariubs/go-producthunt"
	"github.com/dariubs/huntline/app/db"
	"github.com/dariubs/huntline/app/model"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var dbs *gorm.DB

// getYesterday returns yesterday's date as a formatted string in PST.
func getYesterday() string {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	return yesterday.Format("2006-01-02")
}

// runAtScheduledTime schedules a given task to run at a specified hour and minute (PST).
func runAtScheduledTime(task func(), hour, minute int) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

	for {
		now := time.Now().In(loc)
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}
		log.Printf("Next run scheduled at: %s (PST)", nextRun)
		time.Sleep(time.Until(nextRun))
		task()
	}
}

// runTaskForDate executes the product fetching and persistence task for a given date.
func runTaskForDate(client producthunt.ProductHunt, date string) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	parsedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		log.Fatalf("Error parsing date %s: %v", date, err)
	}

	products, err := client.GetProductsByRankByDate(date, 5)
	if err != nil {
		log.Fatalf("Error fetching products for date %s: %v", date, err)
	}

	fmt.Printf("Top Products on %s:\n", date)
	for i, product := range products {
		fmt.Printf("ID: %s\nName: %s\nTagline: %s\nWebsite: %s\nRank: %d\n\n",
			product.ID, product.Name, product.Tagline, product.Website, i+1)

		pdc := model.Product{
			Name:    product.Name,
			Tagline: product.Tagline,
			URL:     product.Website,
			Rank:    uint(i + 1),
			Logo:    product.Thumbnail,
			Date:    parsedDate,
		}

		err = pdc.Save(dbs)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	// Define command-line flags.
	runNow := flag.Bool("run-now", true, "Run the task immediately before starting the scheduler")
	dateParam := flag.String("date", "", "Date in format YYYY-MM-DD to fetch data (overrides default 'yesterday')")
	repeatable := flag.Bool("repeat", false, "Set task to run repeatedly according to the schedule (default true)")
	schedule := flag.String("schedule", "00:30", "Schedule time in 24hr format (HH:MM) when the task should run (default 00:30)")
	historical := flag.Bool("historical", false, "If set, run the task for every day from 2013-11-24 to the present day")
	flag.Parse()

	// Validate the date flag if provided.
	if *dateParam != "" {
		if _, err := time.Parse("2006-01-02", *dateParam); err != nil {
			log.Fatalf("Invalid date format for -date flag. Expected YYYY-MM-DD: %v", err)
		}
	}

	parts := strings.Split(*schedule, ":")
	if len(parts) != 2 {
		log.Fatalf("Invalid schedule time format: expected HH:MM")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("Invalid hour in schedule time: %v", err)
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatalf("Invalid minute in schedule time: %v", err)
	}

	// Load environment variables.
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection.
	dbs, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("PH_API_KEY")
	client := producthunt.ProductHunt{APIKey: apiKey}

	// If the historical flag is set, execute the task for each day from 2013-11-24 to today.
	if *historical {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			log.Fatalf("Failed to load timezone: %v", err)
		}
		startDateStr := "2016-07-29"
		startDate, err := time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			log.Fatalf("Error parsing start date: %v", err)
		}
		// Define the end date as today in the specified time zone.
		endDate := time.Now().In(loc)

		// Iterate day by day.
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			log.Printf("Processing date: %s", dateStr)
			runTaskForDate(client, dateStr)

			time.Sleep(20 * time.Second)
		}
		return
	}

	// Define the task function to run for a specific date.
	task := func() {
		var date string
		if *dateParam != "" {
			date = *dateParam
		} else {
			date = getYesterday()
		}
		runTaskForDate(client, date)
	}

	// Execute the task in either scheduled or single-run mode.
	if *repeatable {
		if *runNow {
			task()
		}
		runAtScheduledTime(task, hour, minute)
	} else {
		task()
	}
}
