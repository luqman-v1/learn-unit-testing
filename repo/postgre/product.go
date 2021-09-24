package postgre

import (
	"context"
	"learn-unit-testing/entity"
)

type IProduct interface {
	Detail(context.Context, string) (entity.Product, error)
}

type Product struct{}

func (p Product) Detail(ctx context.Context, s string) (entity.Product, error) {
	panic("implement me")
}

func NewProduct() *Product {
	return &Product{}
}
