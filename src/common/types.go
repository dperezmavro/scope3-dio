package common

import (
	"errors"
	"fmt"
)

// PropertyQuery is the format that the scope 3 api expects for this endpoint
type PropertyQuery struct {
	Channel     string `json:"channel"`
	Country     string `json:"country"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"inventoryId"`
	UtcDateTime string `json:"utcDatetime"`
	Weight      int    `json:"weight"`
}

func (p PropertyQuery) IndexName() string {
	return fmt.Sprintf("%s-%s", p.UtcDateTime, p.InventoryID)
}

func (p *PropertyQuery) Validate() error {
	if p.InventoryID == "" {
		return errors.New("empty value for InventoryId")
	}

	if p.UtcDateTime == "" {
		return errors.New("empty value for UtcDateTime")
	}

	if p.Impressions == 0 {
		p.Impressions = defaultImpressions
	}

	if p.Weight >= defaultWeight {
		p.Weight = 1 // set to minimum cost
	} else if p.Weight >= 0 && p.Weight < defaultWeight {
		p.Weight = defaultWeight - p.Weight // set to cost relative to default cost
	} else {
		p.Weight = defaultWeight // set to default cost
	}

	return nil
}

// MeasureAPIRequest is the incoming request struct
type MeasureAPIRequest struct {
	Rows []PropertyQuery `json:"rows"`
}

// PropertyResponse is a custom type defined for moving data around here
type PropertyResponse struct {
	PropertyName   string `json:"propertyName,omitempty"`
	UtcDateTime    string `json:"utcDateTime,omitempty"`
	TotalEmissions string `json:"totalEmissions"`
	Weight         int    `json:"weight,omitempty"`
}

func (p PropertyResponse) IndexName() string {
	return fmt.Sprintf("%s-%s", p.UtcDateTime, p.PropertyName)
}
