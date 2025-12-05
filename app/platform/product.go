package platform

import "time"

// Product represents a product from a launch platform
type Product struct {
	Name        string
	URL         string
	Tagline     string
	Description string
	Rank        uint
	Logo        string
	Date        time.Time
	Platform    string
}

