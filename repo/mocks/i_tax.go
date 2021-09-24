// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "learn-unit-testing/entity"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// ITax is an autogenerated mock type for the ITax type
type ITax struct {
	mock.Mock
}

// Insert provides a mock function with given fields: _a0, _a1, _a2
func (_m *ITax) Insert(_a0 context.Context, _a1 *sqlx.Tx, _a2 entity.Tax) (entity.Tax, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 entity.Tax
	if rf, ok := ret.Get(0).(func(context.Context, *sqlx.Tx, entity.Tax) entity.Tax); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(entity.Tax)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sqlx.Tx, entity.Tax) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
