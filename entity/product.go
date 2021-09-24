package entity

//go:generate easyjson -all product.go
type Product struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Name  string  `json:"name"`
}
