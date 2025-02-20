package huntline

import (
	"github.com/dariubs/huntline/app/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IndexHandler(db *gorm.DB, gd types.General) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"gd": gd})
	}
}
