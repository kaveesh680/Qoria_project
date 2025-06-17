package entity

type MonthlySalesVolume struct {
	Month      string `json:"month"`
	TotalSales int    `json:"total_sales"`
}

type MonthlySalesRequest struct {
	Limit int `json:"limit"`
}
