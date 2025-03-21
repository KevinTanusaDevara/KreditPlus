// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	domain "kreditplus/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// TransactionUsecase is an autogenerated mock type for the TransactionUsecase type
type TransactionUsecase struct {
	mock.Mock
}

// CreateTransactionWithLimitUpdate provides a mock function with given fields: userID, customer, input
func (_m *TransactionUsecase) CreateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, input domain.TransactionInput) error {
	ret := _m.Called(userID, customer, input)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransactionWithLimitUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, *domain.Customer, domain.TransactionInput) error); ok {
		r0 = rf(userID, customer, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTransactionWithLimitUpdate provides a mock function with given fields: userID, transaction
func (_m *TransactionUsecase) DeleteTransactionWithLimitUpdate(userID uint, transaction *domain.Transaction) error {
	ret := _m.Called(userID, transaction)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTransactionWithLimitUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, *domain.Transaction) error); ok {
		r0 = rf(userID, transaction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTransactions provides a mock function with given fields: limit, offset
func (_m *TransactionUsecase) GetAllTransactions(limit int, offset int) ([]domain.Transaction, error) {
	ret := _m.Called(limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetAllTransactions")
	}

	var r0 []domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]domain.Transaction, error)); ok {
		return rf(limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int, int) []domain.Transaction); ok {
		r0 = rf(limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomerByNIK provides a mock function with given fields: nik
func (_m *TransactionUsecase) GetCustomerByNIK(nik string) (*domain.Customer, error) {
	ret := _m.Called(nik)

	if len(ret) == 0 {
		panic("no return value specified for GetCustomerByNIK")
	}

	var r0 *domain.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Customer, error)); ok {
		return rf(nik)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Customer); ok {
		r0 = rf(nik)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(nik)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByID provides a mock function with given fields: id
func (_m *TransactionUsecase) GetTransactionByID(id uint) (*domain.Transaction, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionByID")
	}

	var r0 *domain.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*domain.Transaction, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *domain.Transaction); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTransactionWithLimitUpdate provides a mock function with given fields: userID, customer, transaction, input
func (_m *TransactionUsecase) UpdateTransactionWithLimitUpdate(userID uint, customer *domain.Customer, transaction *domain.Transaction, input domain.TransactionInput) error {
	ret := _m.Called(userID, customer, transaction, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTransactionWithLimitUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, *domain.Customer, *domain.Transaction, domain.TransactionInput) error); ok {
		r0 = rf(userID, customer, transaction, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactionUsecase creates a new instance of TransactionUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionUsecase {
	mock := &TransactionUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
