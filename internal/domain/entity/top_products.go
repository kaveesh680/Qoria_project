package entity

type TopProduct struct {
	ProductId      string `json:"product_id"`
	ProductName    string `json:"productName"`
	PurchaseCount  int    `json:"purchaseCount"`
	AvailableStock int    `json:"availableStock"`
}
