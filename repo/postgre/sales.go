package postgre

import (
	"context"
	"learn-unit-testing/entity"

	"github.com/jmoiron/sqlx"
)

type ISales interface {
	Insert(context.Context, *sqlx.Tx, entity.SalesModel) (entity.SalesModel, error)
}

type Sales struct{}

func (s Sales) Insert(ctx context.Context, tx *sqlx.Tx, model entity.SalesModel) (entity.SalesModel, error) {
	panic("implement me")
}

func NewSales() *Sales {
	return &Sales{}
}
