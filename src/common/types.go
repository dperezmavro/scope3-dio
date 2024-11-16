package common

// PropertyQuery is the format that the scope 3 api expects for this endpoint
type PropertyQuery struct {
	Channel     string `json:"channel"`
	Country     string `json:"country"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"inventoryId"`
	UtcDateTime string `json:"utcDatetime"`
}

// MeasureAPIRequest is the incoming request struct
type MeasureAPIRequest struct {
	Rows []PropertyQuery `json:"rows"`
}

// PropertyResponse is a custom type defined for moving data around here
type PropertyResponse struct {
	InventoryID string `json:"inventoryId"`
	UtcDateTime string `json:"utcDatetime"`
	Body        string
}
