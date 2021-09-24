package entity

type SalesModel struct {
	ID            string  `json:"id"`
	TotalBill     float64 `json:"total_bill"`
	TotalDiscount float64 `json:"total_discount"`
	TotalTax      float64 `json:"total_tax"`
	Quantity      float64 `json:"quantity"`
	ProductID     string  `json:"product"`
	Price         float64 `json:"price"`
}

type Tax struct {
	ID      string  `json:"id"`
	SalesID string  `json:"sales_id"`
	Amount  float64 `json:"amount"`
}

type Discount struct {
	ID      string  `json:"id"`
	SalesID string  `json:"sales_id"`
	Amount  float64 `json:"amount"`
}
