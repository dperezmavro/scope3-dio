package scope3

import "net/http"

// HTTPClient interface for swaping out clients (useful in testing)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MeasureAPIRequest is the incoming request struct
type MeasureAPIRequest struct {
	Rows []RowItem `json:"rows"`
}

// RowItem is an item in the Rows slice of MeasureAPIRequest
type RowItem struct {
	Channel     string `json:"channel,omitempty"`
	Country     string `json:"country,omitempty"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"InventoryId"`
	UtcDateTime string `json:"utcDatetime"`
}

// MeasureAPIResponse represents the entire JSON response structure.
type MeasureAPIResponse struct {
	Rows []Row `json:"rows"`
}

// Row represents an individual row of emissions data.
type Row struct {
	TotalEmissions float64 `json:"totalEmissions"`
	InventoryID    string  `json:"inventoryId"`
}
