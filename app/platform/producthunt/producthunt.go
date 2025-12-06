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

	// Ensure parsedDate is normalized to midnight in PST to avoid timezone conversion issues
	// This ensures the date stays as the requested date when saved to the database
	normalizedDate := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)

	products := make([]platform.Product, len(phProducts))
	for i, phProduct := range phProducts {
		// Generate Google favicon URL from the product's website URL
		// This ensures we get the actual product logo, not ProductHunt's thumbnail
		logoURL := ""
		if phProduct.Website != "" {
			logoURL = "https://www.google.com/s2/favicons?domain=" + phProduct.Website + "&sz=64"
		}

		products[i] = platform.Product{
			Name:        phProduct.Name,
			Tagline:     phProduct.Tagline,
			URL:         phProduct.Website,
			Rank:        uint(i + 1),
			Logo:        logoURL, // Use Google favicon service for correct product logos
			Date:        normalizedDate, // Use normalized date to ensure correct date is saved
			Platform:    PlatformName,
			Description: "", // ProductHunt API doesn't provide description in the current struct
		}
	}

	return products, nil
}

