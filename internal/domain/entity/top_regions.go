package entity

type RegionRevenue struct {
	Region       string  `json:"region"`
	TotalRevenue float64 `json:"total_revenue"`
	ItemsSold    int     `json:"items_sold"`
}

type TopRegionsRequest struct {
	Limit int `json:"limit"`
}
