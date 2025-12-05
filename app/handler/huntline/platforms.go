package huntline

import (
	"net/http"

	"github.com/dariubs/huntline/app/model"
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PlatformsHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all distinct platforms with statistics
		type PlatformStats struct {
			Platform     string
			ProductCount int64
			EarliestDate string
			LatestDate   string
			DateCount    int64
		}

		var platforms []string
		db.Model(&model.Product{}).Distinct("platform").Pluck("platform", &platforms)

		var platformStatsList []PlatformStats
		for _, platform := range platforms {
			var productCount int64
			var earliestDate, latestDate string
			var dateCount int64

			db.Model(&model.Product{}).Where("platform = ?", platform).Count(&productCount)

			var dates []struct {
				Date string
			}
			db.Model(&model.Product{}).Select("DISTINCT date").Where("platform = ?", platform).
				Order("date ASC").Limit(1).Scan(&dates)
			if len(dates) > 0 {
				earliestDate = dates[0].Date
			}

			db.Model(&model.Product{}).Select("DISTINCT date").Where("platform = ?", platform).
				Order("date DESC").Limit(1).Scan(&dates)
			if len(dates) > 0 {
				latestDate = dates[0].Date
			}

			db.Model(&model.Product{}).Select("COUNT(DISTINCT date)").Where("platform = ?", platform).
				Scan(&dateCount)

			platformStatsList = append(platformStatsList, PlatformStats{
				Platform:     platform,
				ProductCount: productCount,
				EarliestDate: earliestDate,
				LatestDate:   latestDate,
				DateCount:    dateCount,
			})
		}

		c.HTML(http.StatusOK, "platforms.html", gin.H{
			"gd":          gd,
			"platforms":   platformStatsList,
			"currentPage": "platforms",
		})
	}
}
