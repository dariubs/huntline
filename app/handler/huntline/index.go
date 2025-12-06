package huntline

import (
	"time"

	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlatformProducts groups products by platform
type PlatformProducts struct {
	Platform string
	Products []model.Product
}

func IndexHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get dates in San Francisco timezone (Pacific Time)
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			c.String(500, "Failed to load timezone")
			return
		}
		now := time.Now().In(loc)
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		yesterday := today.AddDate(0, 0, -1)
		
		// Get date parameter (format: YYYY-MM-DD) or use today as default
		dateParam := c.DefaultQuery("date", "")
		var startDate, endDate time.Time
		
		if dateParam != "" {
			parsedDate, err := time.ParseInLocation("2006-01-02", dateParam, loc)
			if err == nil {
				startDate = parsedDate
				endDate = parsedDate
			} else {
				// On error, default to today
				startDate = today
				endDate = today
			}
		} else {
			// Default: show today's products
			// Redirect to include date parameter in URL for clarity
			todayStr := today.Format("2006-01-02")
			c.Redirect(302, "/?date="+todayStr)
			return
		}
		
		// Query products grouped by platform for the date range
		// Format dates as strings for PostgreSQL DATE field comparison
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

		// Group products by platform
		platformMap := make(map[string][]model.Product)
		for _, product := range allProducts {
			platformMap[product.Platform] = append(platformMap[product.Platform], product)
		}

		// Convert to slice for template
		var platformProducts []PlatformProducts
		for platform, products := range platformMap {
			platformProducts = append(platformProducts, PlatformProducts{
				Platform: platform,
				Products: products,
			})
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
		for _, pp := range platformProducts {
			// Group by date
			dateMap := make(map[string][]model.Product)
			for _, product := range pp.Products {
				dateStr := product.Date.Format("2006-01-02")
				dateMap[dateStr] = append(dateMap[dateStr], product)
			}
			
			var dateGroups []DateGroup
			for dateStr, products := range dateMap {
				parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, loc)
				if err != nil {
					// Skip invalid dates
					continue
				}
				dateGroups = append(dateGroups, DateGroup{
					Date:     parsedDate,
					DateStr:  dateStr,
					Products: products,
				})
			}
			
			platformDataList = append(platformDataList, PlatformData{
				Platform:   pp.Platform,
				DateGroups: dateGroups,
			})
		}

		// Determine if showing today or a specific date
		isToday := endDate.Format("2006-01-02") == today.Format("2006-01-02")
		
		// Format date for display: "2 January 2006"
		dateFormatted := endDate.Format("2 January 2006")
		
		c.HTML(200, "index.html", gin.H{
			"gd":               gd,
			"platforms":        platformDataList,
			"todayStr":         today.Format("2006-01-02"),
			"yesterdayStr":     yesterday.Format("2006-01-02"),
			"currentDate":      endDate.Format("2006-01-02"),
			"currentDateFormatted": dateFormatted,
			"isToday":          isToday,
			"prevDay":          prevDay.Format("2006-01-02"),
			"nextDay":          nextDay.Format("2006-01-02"),
			"nextDayInFuture":  nextDayInFuture,
			"startDate":        startDate.Format("2006-01-02"),
			"endDate":          endDate.Format("2006-01-02"),
		})
	}
}
