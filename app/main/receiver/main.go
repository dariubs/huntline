package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dariubs/huntline/app/db"
	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/platform"
	"github.com/dariubs/huntline/app/platform/producthunt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var dbs *gorm.DB

// getYesterday returns yesterday's date as a formatted string in San Francisco timezone (Pacific Time).
func getYesterday() string {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	return yesterday.Format("2006-01-02")
}

// runAtScheduledTime schedules a given task to run at a specified hour and minute (San Francisco timezone - Pacific Time).
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
		log.Printf("Next run scheduled at: %s (San Francisco/Pacific Time)", nextRun)
		time.Sleep(time.Until(nextRun))
		task()
	}
}

// runTaskForDate executes the product fetching and persistence task for a given date and platform.
func runTaskForDate(platformClient platform.LaunchPlatform, date string) {
	products, err := platformClient.GetTopProducts(date, 10)
	if err != nil {
		log.Fatalf("Error fetching products for platform %s on date %s: %v", platformClient.GetName(), date, err)
	}

	fmt.Printf("Top Products from %s on %s:\n", platformClient.GetName(), date)
	for _, product := range products {
		fmt.Printf("Name: %s\nTagline: %s\nWebsite: %s\nRank: %d\nPlatform: %s\n\n",
			product.Name, product.Tagline, product.URL, product.Rank, product.Platform)

		pdc := model.Product{
			Name:        product.Name,
			Tagline:     product.Tagline,
			URL:         product.URL,
			Rank:        product.Rank,
			Logo:        product.Logo,
			Date:        product.Date,
			Platform:    product.Platform,
			Description: product.Description,
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
	historical := flag.Bool("historical", false, "If set, run the task for every day from 2016-07-29 to the present day")
	lastMonth := flag.Bool("last-month", false, "If set, run the task for every day in the previous month")
	platformParam := flag.String("platform", "producthunt", "Platform to fetch products from (default: producthunt)")
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

	// Initialize platform client based on platform parameter
	var platformClient platform.LaunchPlatform
	switch *platformParam {
	case "producthunt":
		apiKey := os.Getenv("PH_API_KEY")
		if apiKey == "" {
			log.Fatal("PH_API_KEY environment variable is required for ProductHunt platform")
		}
		platformClient = producthunt.NewProductHuntPlatform(apiKey)
	default:
		log.Fatalf("Unsupported platform: %s. Supported platforms: producthunt", *platformParam)
	}

	// If the historical flag is set, execute the task for each day from 2016-07-29 to today.
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
			log.Printf("Processing date: %s for platform: %s", dateStr, platformClient.GetName())
			runTaskForDate(platformClient, dateStr)

			time.Sleep(20 * time.Second)
		}
		return
	}

	// If the last-month flag is set, execute the task for each day in the previous month.
	if *lastMonth {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			log.Fatalf("Failed to load timezone: %v", err)
		}
		now := time.Now().In(loc)
		
		// Calculate first day of current month
		firstOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		
		// Calculate first day of last month
		firstOfLastMonth := firstOfCurrentMonth.AddDate(0, -1, 0)
		
		// Calculate last day of last month (first day of current month minus 1 day)
		lastOfLastMonth := firstOfCurrentMonth.AddDate(0, 0, -1)

		log.Printf("Updating last month (%s to %s) for platform: %s", 
			firstOfLastMonth.Format("2006-01-02"), 
			lastOfLastMonth.Format("2006-01-02"),
			platformClient.GetName())

		// Iterate day by day through last month.
		for d := firstOfLastMonth; !d.After(lastOfLastMonth); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			log.Printf("Processing date: %s for platform: %s", dateStr, platformClient.GetName())
			runTaskForDate(platformClient, dateStr)

			// Small delay to avoid rate limiting
			time.Sleep(5 * time.Second)
		}
		
		log.Printf("Finished updating last month's data for platform: %s", platformClient.GetName())
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
		runTaskForDate(platformClient, date)
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
