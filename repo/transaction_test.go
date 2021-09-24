package repo

import (
	"context"
	"errors"
	"learn-unit-testing/entity"
	"learn-unit-testing/repo/mocks"
	"learn-unit-testing/repo/postgre"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type TrxTestSuite struct {
	suite.Suite
	sales    *mocks.ISales
	tax      *mocks.ITax
	discount *mocks.IDiscount
	product  *mocks.IProduct
	redis    *mocks.IRedis
	act      *Transaction
}

func (s *TrxTestSuite) SetupTest() {
	s.sales = new(mocks.ISales)
	s.tax = new(mocks.ITax)
	s.discount = new(mocks.IDiscount)
	s.product = new(mocks.IProduct)
	s.redis = new(mocks.IRedis)

	s.act = NewTransaction()
	s.act.sales = s.sales
	s.act.tax = s.tax
	s.act.discount = s.discount
	s.act.product = s.product
	s.act.redis = s.redis

}

func TestTrxTestSuite(t *testing.T) {
	suite.Run(t, new(TrxTestSuite))
}

func (s *TrxTestSuite) AfterTest(_, _ string) {
	s.sales.AssertExpectations(s.T())
	s.tax.AssertExpectations(s.T())
	s.discount.AssertExpectations(s.T())
	s.product.AssertExpectations(s.T())
	s.redis.AssertExpectations(s.T())
}

func (s *TrxTestSuite) TestTransaction_CalculateTotalTax() {
	assert.Equal(s.T(), float64(100), s.act.CalculateTotalTax(1000, 10))
}
func (s *TrxTestSuite) TestTransaction_CalculateTotalBill() {
	assert.Equal(s.T(), float64(10000), s.act.CalculateTotalBill(1000, 10))
}

func (s *TrxTestSuite) TestTransaction_CalculateTotalDiscount() {
	assert.Equal(s.T(), float64(9000), s.act.CalculateTotalDiscount(10000, 1000))
}

func (s *TrxTestSuite) TestTransaction_Order() {
	request := entity.Request{
		Quantity:  10,
		Discount:  1000,
		Tax:       10,
		ProductID: "1",
	}
	s.redis.On("Get", ctx, mock.Anything).Return(nil, errors.New("not found"))

	s.product.On("Detail", ctx, request.ProductID).Return(entity.Product{
		ID:    "1",
		Price: 1000,
		Name:  "cola-cola",
	}, nil)

	s.redis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockDB, sqlMock, _ := sqlmock.New()
	sqlMock.ExpectBegin()
	defer mockDB.Close()

	postgre.DB = sqlx.NewDb(mockDB, "sqlmock")

	salesDummy := entity.SalesModel{
		TotalBill:     10000,
		TotalDiscount: 9000,
		TotalTax:      1000,
		Quantity:      10,
		ProductID:     "1",
		Price:         1000,
	}
	s.sales.On("Insert", ctx, mock.Anything, salesDummy).Return(salesDummy, nil)

	s.discount.On("Insert", ctx, mock.Anything, mock.Anything).Return(entity.Discount{
		SalesID: "1",
		Amount:  request.Discount,
	}, nil)

	s.tax.On("Insert", ctx, mock.Anything, mock.Anything).Return(entity.Tax{
		SalesID: "1",
		Amount:  request.Tax,
	}, nil)
	sqlMock.ExpectCommit()
	resp, err := s.act.Order(ctx, request)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), float64(10000), resp.TotalBill)
	assert.Equal(s.T(), float64(9000), resp.TotalDiscount)
	assert.Equal(s.T(), float64(1000), resp.TotalTax)
}

func (s *TrxTestSuite) TestTransaction_FindProductOnRedis() {
	productID := "1"
	responseDummy := entity.Product{
		ID:    "1",
		Price: 1000,
		Name:  "cola-cola",
	}
	b, _ := responseDummy.MarshalJSON()

	s.redis.On("Get", ctx, mock.Anything).Return(b, nil)
	product, err := s.act.FindProduct(ctx, productID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), responseDummy, product)
}

func (s *TrxTestSuite) TestTransaction_FindProductNotInRedis() {
	productID := "1"
	responseDummy := entity.Product{
		ID:    "1",
		Price: 1000,
		Name:  "cola-cola",
	}
	s.redis.On("Get", ctx, mock.Anything).Return(nil, errors.New("not found"))
	s.product.On("Detail", ctx, productID).Return(responseDummy, nil)
	s.redis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	product, err := s.act.FindProduct(ctx, productID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), responseDummy, product)
}
