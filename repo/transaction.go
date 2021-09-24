package repo

import (
	"context"
	"learn-unit-testing/entity"
	"learn-unit-testing/repo/postgre"
	"learn-unit-testing/repo/redis"
	"time"
)

type Transaction struct {
	sales    postgre.ISales
	tax      postgre.ITax
	discount postgre.IDiscount
	product  postgre.IProduct
	redis    redis.IRedis
}

func NewTransaction() *Transaction {
	return &Transaction{
		sales:    postgre.NewSales(),
		tax:      postgre.NewTax(),
		discount: postgre.NewDiscount(),
		product:  postgre.NewProduct(),
		redis:    redis.NewRedisDB(),
	}
}

func (t *Transaction) CalculateTotalBill(price, quantity float64) float64 {
	return price * quantity
}

func (t *Transaction) CalculateTotalTax(totalBill, tax float64) float64 {
	return (tax / 100) * totalBill
}

func (t *Transaction) CalculateTotalDiscount(totalBill, discount float64) float64 {
	return totalBill - discount
}

func (t *Transaction) keyRedis(productID string) string {
	return "product:detail:" + productID
}

func (t *Transaction) FindProduct(ctx context.Context, productID string) (entity.Product, error) {
	var product entity.Product
	keyRedis := t.keyRedis(productID)
	get, err := t.redis.Get(ctx, keyRedis)
	if err == nil {
		err := product.UnmarshalJSON(get)
		if err != nil {
			return entity.Product{}, err
		}
	} else {
		detail, err := t.product.Detail(ctx, productID)
		if err != nil {
			return entity.Product{}, err
		}
		b, err := detail.MarshalJSON()
		if err != nil {
			return entity.Product{}, err
		}
		_ = t.redis.Set(ctx, keyRedis, b, 1*time.Minute)
		product = detail
	}
	return product, nil
}

func (t *Transaction) Order(ctx context.Context, request entity.Request) (entity.SalesModel, error) {
	product, err := t.FindProduct(ctx, request.ProductID)
	if err != nil {
		return entity.SalesModel{}, err
	}
	tx, err := postgre.DB.Beginx()
	if err != nil {
		return entity.SalesModel{}, err
	}
	totalBill := t.CalculateTotalBill(product.Price, request.Quantity)

	respSales, err := t.sales.Insert(ctx, tx, entity.SalesModel{
		TotalBill:     totalBill,
		TotalDiscount: t.CalculateTotalDiscount(totalBill, request.Discount),
		TotalTax:      t.CalculateTotalTax(totalBill, request.Tax),
		Quantity:      request.Quantity,
		Price:         product.Price,
		ProductID:     request.ProductID,
	})

	if err != nil {
		return entity.SalesModel{}, err
	}

	_, err = t.discount.Insert(ctx, tx, entity.Discount{
		SalesID: respSales.ID,
		Amount:  request.Discount,
	})

	if err != nil {
		return entity.SalesModel{}, err
	}

	_, err = t.tax.Insert(ctx, tx, entity.Tax{
		SalesID: respSales.ID,
		Amount:  request.Tax,
	})

	if err != nil {
		return entity.SalesModel{}, err
	}
	err = tx.Commit()
	if err != nil {
		return entity.SalesModel{}, err
	}
	return respSales, nil
}
