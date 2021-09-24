package entity

type Request struct {
	Quantity  float64 `json:"quantity"`
	Discount  float64 `json:"discount"`
	Tax       float64 `json:"tax"`
	ProductID string  `json:"product_id"`
}
