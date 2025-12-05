package huntline

import (
	"net/http"
	"time"

	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ArchiveHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		platform := c.Query("platform")
		if platform == "" {
			platform = "producthunt" // Default to producthunt
		}

		// Get dates in San Francisco timezone (Pacific Time)
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			c.String(500, "Failed to load timezone")
			return
		}
		now := time.Now().In(loc)

		// Get month/year from query params or use current month
		monthStr := c.DefaultQuery("month", "")
		var selectedMonth time.Time

		if monthStr != "" {
			parsedMonth, err := time.Parse("2006-01", monthStr)
			if err == nil {
				selectedMonth = parsedMonth
			} else {
				selectedMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
			}
		} else {
			selectedMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		}

		// Calculate month boundaries
		monthStart := time.Date(selectedMonth.Year(), selectedMonth.Month(), 1, 0, 0, 0, 0, loc)
		monthEnd := monthStart.AddDate(0, 1, 0).AddDate(0, 0, -1)

		// Get all distinct dates for the platform in the selected month, ordered by date descending
		var dates []struct {
			Date time.Time
		}

		query := db.Model(&model.Product{}).Select("DISTINCT date").
			Where("date >= ? AND date <= ?", monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).
			Order("date DESC")
		if platform != "all" {
			query = query.Where("platform = ?", platform)
		}
		query.Scan(&dates)

		// For each date, get products
		type DateGroup struct {
			Date     time.Time
			DateStr  string
			Products []model.Product
		}

		var dateGroups []DateGroup
		for _, d := range dates {
			var products []model.Product
			productQuery := db.Where("date = ?", d.Date.Format("2006-01-02")).Order("rank ASC")
			if platform != "all" {
				productQuery = productQuery.Where("platform = ?", platform)
			}
			productQuery.Find(&products)

			if len(products) > 0 {
				dateGroups = append(dateGroups, DateGroup{
					Date:     d.Date,
					DateStr:  d.Date.Format("2006-01-02"),
					Products: products,
				})
			}
		}

		// Calculate previous and next month for navigation
		prevMonth := monthStart.AddDate(0, -1, 0)
		nextMonth := monthStart.AddDate(0, 1, 0)

		// Check if next month is in the future
		nextMonthInFuture := nextMonth.After(now)

		c.HTML(http.StatusOK, "archive.html", gin.H{
			"gd":                gd,
			"dateGroups":        dateGroups,
			"platform":          platform,
			"currentMonth":      monthStart.Format("2006-01"),
			"currentMonthName":  monthStart.Format("January 2006"),
			"prevMonth":         prevMonth.Format("2006-01"),
			"nextMonth":         nextMonth.Format("2006-01"),
			"nextMonthInFuture": nextMonthInFuture,
			"currentPage":       "archive",
		})
	}
}
