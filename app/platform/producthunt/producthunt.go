package producthunt

import (
	"github.com/dariubs/go-producthunt"
	"github.com/dariubs/huntline/app/platform"
	"time"
)

const PlatformName = "producthunt"

type ProductHuntPlatform struct {
	client producthunt.ProductHunt
}

// NewProductHuntPlatform creates a new ProductHunt platform instance
func NewProductHuntPlatform(apiKey string) *ProductHuntPlatform {
	return &ProductHuntPlatform{
		client: producthunt.ProductHunt{APIKey: apiKey},
	}
}

// GetName returns the platform name
func (p *ProductHuntPlatform) GetName() string {
	return PlatformName
}

// GetTopProducts fetches top products from ProductHunt for a given date
func (p *ProductHuntPlatform) GetTopProducts(date string, limit int) ([]platform.Product, error) {
	// Use San Francisco timezone (Pacific Time)
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}
	
	parsedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, err
	}

	phProducts, err := p.client.GetProductsByRankByDate(date)
	if err != nil {
		return nil, err
	}

	// Limit the number of products if needed
	if limit > 0 && limit < len(phProducts) {
		phProducts = phProducts[:limit]
	}

	products := make([]platform.Product, len(phProducts))
	for i, phProduct := range phProducts {
		products[i] = platform.Product{
			Name:        phProduct.Name,
			Tagline:     phProduct.Tagline,
			URL:         phProduct.Website,
			Rank:        uint(i + 1),
			Logo:        phProduct.Thumbnail,
			Date:        parsedDate,
			Platform:    PlatformName,
			Description: "", // ProductHunt API doesn't provide description in the current struct
		}
	}

	return products, nil
}

