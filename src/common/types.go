package common

import "fmt"

// PropertyQuery is the format that the scope 3 api expects for this endpoint
type PropertyQuery struct {
	Channel     string `json:"channel"`
	Country     string `json:"country"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"inventoryId"`
	UtcDateTime string `json:"utcDatetime"`
	Weight      int64  `json:"weight"`
}

func (p PropertyQuery) IndexName() string {
	return fmt.Sprintf("%s-%s", p.UtcDateTime, p.InventoryID)
}

// MeasureAPIRequest is the incoming request struct
type MeasureAPIRequest struct {
	Rows []PropertyQuery `json:"rows"`
}

// PropertyResponse is a custom type defined for moving data around here
type PropertyResponse struct {
	InventoryID string
	UtcDateTime string
	Body        string
	Weight      int64
}

func (p PropertyResponse) IndexName() string {
	return fmt.Sprintf("%s-%s", p.UtcDateTime, p.InventoryID)
}
