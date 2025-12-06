package huntline

import (
	"time"

	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IndexAPIHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get dates in San Francisco timezone (Pacific Time)
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load timezone"})
			return
		}
		now := time.Now().In(loc)
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		yesterday := today.AddDate(0, 0, -1)
		
		// Get date parameter (format: YYYY-MM-DD) or use today
		dateParam := c.DefaultQuery("date", "")
		var startDate, endDate time.Time
		
		if dateParam != "" {
			parsedDate, err := time.ParseInLocation("2006-01-02", dateParam, loc)
			if err == nil {
				startDate = parsedDate
				endDate = parsedDate
			} else {
				startDate = today
				endDate = today
			}
		} else {
			// Default: show today's products
			startDate = today
			endDate = today
		}
		
		// Query products grouped by platform for the date range
		startDateStr := startDate.Format("2006-01-02")
		endDateStr := endDate.Format("2006-01-02")
		
		var allProducts []model.Product
		db.Where("date >= ? AND date <= ?", startDateStr, endDateStr).
			Order("platform ASC, date DESC, rank ASC").
			Find(&allProducts)
		
		// Calculate navigation dates
		prevDay := startDate.AddDate(0, 0, -1)
		nextDay := endDate.AddDate(0, 0, 1)
		nextDayInFuture := nextDay.After(today)
		
		// Determine if showing today
		isToday := endDate.Format("2006-01-02") == today.Format("2006-01-02")

		// Group products by platform
		platformMap := make(map[string][]model.Product)
		for _, product := range allProducts {
			platformMap[product.Platform] = append(platformMap[product.Platform], product)
		}

		// Group products by date within each platform
		type DateGroup struct {
			Date     time.Time
			DateStr  string
			Products []model.Product
		}
		
		type PlatformData struct {
			Platform   string
			DateGroups []DateGroup
		}
		
		var platformDataList []PlatformData
		for platform, products := range platformMap {
			// Group by date
			dateMap := make(map[string][]model.Product)
			for _, product := range products {
				dateStr := product.Date.Format("2006-01-02")
				dateMap[dateStr] = append(dateMap[dateStr], product)
			}
			
			var dateGroups []DateGroup
			for dateStr, products := range dateMap {
				parsedDate, _ := time.Parse("2006-01-02", dateStr)
				dateGroups = append(dateGroups, DateGroup{
					Date:     parsedDate,
					DateStr:  dateStr,
					Products: products,
				})
			}
			
			platformDataList = append(platformDataList, PlatformData{
				Platform:   platform,
				DateGroups: dateGroups,
			})
		}

		// Return JSON response
		c.JSON(200, gin.H{
			"platforms":        platformDataList,
			"todayStr":         today.Format("2006-01-02"),
			"yesterdayStr":     yesterday.Format("2006-01-02"),
			"currentDate":      endDate.Format("2006-01-02"),
			"currentDateObj":   endDate.Format("2006-01-02"), // For JavaScript parsing
			"isToday":          isToday,
			"prevDay":          prevDay.Format("2006-01-02"),
			"nextDay":          nextDay.Format("2006-01-02"),
			"nextDayInFuture":  nextDayInFuture,
			"startDate":        startDate.Format("2006-01-02"),
			"endDate":          endDate.Format("2006-01-02"),
		})
	}
}

