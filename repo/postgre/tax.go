package postgre

import (
	"context"
	"learn-unit-testing/entity"

	"github.com/jmoiron/sqlx"
)

type ITax interface {
	Insert(context.Context, *sqlx.Tx, entity.Tax) (entity.Tax, error)
}

type Tax struct{}

func (t Tax) Insert(ctx context.Context, tx *sqlx.Tx, tax entity.Tax) (entity.Tax, error) {
	panic("implement me")
}

func NewTax() *Tax {
	return &Tax{}
}
