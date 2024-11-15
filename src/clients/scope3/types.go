package scope3

type PropertyQuery struct {
	Channel     string `json:"channel"`
	Country     string `json:"country"`
	InventoryID string `json:"inventoryId"`
	UtcDateTime string `json:"utcDatetime"`
}

type PropertyResponse struct {
	PropertyQuery
	Body string
}

// MeasureAPIRequest is the incoming request struct
type MeasureAPIRequest struct {
	Rows []RowItem `json:"rows"`
}

// RowItem is an item in the Rows slice of MeasureAPIRequest
type RowItem struct {
	Channel     string `json:"channel"`
	Country     string `json:"country"`
	Impressions int    `json:"impressions"`
	InventoryID string `json:"InventoryId"`
	UtcDateTime string `json:"utcDatetime"`
}

// MeasureAPIResponse represents the entire JSON response structure.
type MeasureAPIResponse struct {
	Coverage                Coverage           `json:"coverage"`
	Policies                []Policy           `json:"policies"`
	RequestID               string             `json:"requestId"`
	TotalEmissions          float64            `json:"totalEmissions"`
	TotalEmissionsBreakdown EmissionsBreakdown `json:"totalEmissionsBreakdown"`
	Rows                    []Row              `json:"rows"`
}

// Coverage contains the coverage information for the emissions data.
type Coverage struct {
	AdFormats        AdFormatMetrics `json:"adFormats"`
	Channels         ChannelMetrics  `json:"channels"`
	MediaOwners      MetricStatus    `json:"mediaOwners"`
	Properties       MetricStatus    `json:"properties"`
	Sellers          MetricStatus    `json:"sellers"`
	TotalImpressions TotalMetric     `json:"totalImpressions"`
	TotalRows        TotalMetric     `json:"totalRows"`
}

// AdFormatMetrics represents the breakdown of ad formats.
type AdFormatMetrics struct {
	Generic        int    `json:"generic"`
	Metric         string `json:"metric"`
	Unknown        int    `json:"unknown"`
	VendorSpecific int    `json:"vendorSpecific"`
}

// ChannelMetrics represents the breakdown of channels.
type ChannelMetrics struct {
	Deprecated int    `json:"deprecated"`
	Metric     string `json:"metric"`
	Modeled    int    `json:"modeled"`
	Unknown    int    `json:"unknown"`
}

// MetricStatus represents the metric status of media owners, properties, sellers.
type MetricStatus struct {
	Metric  string `json:"metric"`
	Modeled int    `json:"modeled"`
	Unknown int    `json:"unknown"`
}

// TotalMetric represents the breakdown of total metrics like impressions and rows.
type TotalMetric struct {
	Metric  string `json:"metric"`
	Modeled int    `json:"modeled"`
	Skipped int    `json:"skipped"`
}

// Policy represents a policy and its compliance status.
type Policy struct {
	Compliant    int    `json:"compliant"`
	Noncompliant int    `json:"noncompliant"`
	Policy       string `json:"policy"`
	PolicyOwner  string `json:"policyOwner"`
}

// EmissionsBreakdown represents the emissions breakdown for different categories.
type EmissionsBreakdown struct {
	Framework string            `json:"framework"`
	Totals    EmissionsCategory `json:"totals"`
}

// EmissionsCategory represents the total emissions for categories like adSelection, creativeDelivery, mediaDistribution.
type EmissionsCategory struct {
	AdSelection       float64 `json:"adSelection"`
	CreativeDelivery  float64 `json:"creativeDelivery"`
	MediaDistribution float64 `json:"mediaDistribution"`
}

// Row represents an individual row of emissions data.
type Row struct {
	EmissionsBreakdown EmissionsBreakdown `json:"emissionsBreakdown"`
	InventoryCoverage  string             `json:"inventoryCoverage"`
	TotalEmissions     float64            `json:"totalEmissions"`
	Internal           InternalData       `json:"internal"`
}

// InternalData represents internal metadata for a particular row.
type InternalData struct {
	CountryRegionGCO2PerKwh int    `json:"countryRegionGCO2PerKwh"`
	CountryRegionCountry    string `json:"countryRegionCountry"`
	Channel                 string `json:"channel"`
	DeviceType              string `json:"deviceType"`
	PropertyID              int    `json:"propertyId"`
	PropertyInventoryType   string `json:"propertyInventoryType"`
	PropertyName            string `json:"propertyName"`
	BenchmarkPercentile     int    `json:"benchmarkPercentile"`
	IsMFA                   bool   `json:"isMFA"`
}
