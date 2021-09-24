package postgre

import (
	"context"
	"learn-unit-testing/entity"

	"github.com/jmoiron/sqlx"
)

type IDiscount interface {
	Insert(context.Context, *sqlx.Tx, entity.Discount) (entity.Discount, error)
}

type Discount struct{}

func (d Discount) Insert(ctx context.Context, tx *sqlx.Tx, discount entity.Discount) (entity.Discount, error) {
	panic("implement me")
}

func NewDiscount() *Discount {
	return &Discount{}
}
