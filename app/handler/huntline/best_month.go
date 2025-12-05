package huntline

import (
	"net/http"
	"time"

	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BestMonthHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get month/year from query params or use current month
		monthStr := c.DefaultQuery("month", "")
		// Get dates in San Francisco timezone (Pacific Time)
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			c.String(500, "Failed to load timezone")
			return
		}
		now := time.Now().In(loc)

		var startDate, endDate time.Time
		if monthStr != "" {
			parsedMonth, err := time.Parse("2006-01", monthStr)
			if err == nil {
				startDate = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, loc)
				endDate = startDate.AddDate(0, 1, 0).AddDate(0, 0, -1)
			} else {
				startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
				endDate = startDate.AddDate(0, 1, 0).AddDate(0, 0, -1)
			}
		} else {
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
			endDate = startDate.AddDate(0, 1, 0).AddDate(0, 0, -1)
		}

		// Get all products for the month, grouped by platform
		var allProducts []model.Product
		db.Where("date >= ? AND date <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
			Order("platform ASC, rank ASC").
			Find(&allProducts)

		// Group by platform and get top products per platform
		type PlatformBest struct {
			Platform string
			Products []model.Product
		}

		platformMap := make(map[string][]model.Product)
		for _, product := range allProducts {
			platformMap[product.Platform] = append(platformMap[product.Platform], product)
		}

		var platformBests []PlatformBest
		for platform, products := range platformMap {
			// Get unique products by name, keeping the one with best rank
			productMap := make(map[string]model.Product)
			for _, p := range products {
				if existing, ok := productMap[p.Name]; !ok || p.Rank < existing.Rank {
					productMap[p.Name] = p
				}
			}

			// Convert to slice and sort by rank
			var uniqueProducts []model.Product
			for _, p := range productMap {
				uniqueProducts = append(uniqueProducts, p)
			}

			// Sort by rank and take top 20
			for i := 0; i < len(uniqueProducts)-1; i++ {
				for j := i + 1; j < len(uniqueProducts); j++ {
					if uniqueProducts[i].Rank > uniqueProducts[j].Rank {
						uniqueProducts[i], uniqueProducts[j] = uniqueProducts[j], uniqueProducts[i]
					}
				}
			}

			if len(uniqueProducts) > 20 {
				uniqueProducts = uniqueProducts[:20]
			}

			platformBests = append(platformBests, PlatformBest{
				Platform: platform,
				Products: uniqueProducts,
			})
		}

		c.HTML(http.StatusOK, "month.html", gin.H{
			"gd":        gd,
			"platforms": platformBests,
			"month":     startDate.Format("January 2006"),
			"monthStr":  startDate.Format("2006-01"),
		})
	}
}
