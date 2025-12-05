package platform

// LaunchPlatform defines the interface that all launch platforms must implement
type LaunchPlatform interface {
	// GetName returns the name/identifier of the platform (e.g., "producthunt", "altern")
	GetName() string
	
	// GetTopProducts fetches the top products for a given date
	// date: date in format YYYY-MM-DD
	// limit: maximum number of products to fetch
	// Returns a slice of Product structs
	GetTopProducts(date string, limit int) ([]Product, error)
}

