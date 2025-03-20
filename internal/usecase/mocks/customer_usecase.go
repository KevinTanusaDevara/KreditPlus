// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	domain "kreditplus/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// CustomerUsecase is an autogenerated mock type for the CustomerUsecase type
type CustomerUsecase struct {
	mock.Mock
}

// CreateCustomer provides a mock function with given fields: input
func (_m *CustomerUsecase) CreateCustomer(input domain.Customer) error {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for CreateCustomer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Customer) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCustomer provides a mock function with given fields: id
func (_m *CustomerUsecase) DeleteCustomer(id uint) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCustomer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllCustomers provides a mock function with given fields: limit, offset
func (_m *CustomerUsecase) GetAllCustomers(limit int, offset int) ([]domain.Customer, error) {
	ret := _m.Called(limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetAllCustomers")
	}

	var r0 []domain.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]domain.Customer, error)); ok {
		return rf(limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int, int) []domain.Customer); ok {
		r0 = rf(limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomerByID provides a mock function with given fields: id
func (_m *CustomerUsecase) GetCustomerByID(id uint) (*domain.Customer, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetCustomerByID")
	}

	var r0 *domain.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*domain.Customer, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *domain.Customer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCustomer provides a mock function with given fields: input
func (_m *CustomerUsecase) UpdateCustomer(input domain.Customer) error {
	ret := _m.Called(input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCustomer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Customer) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCustomerUsecase creates a new instance of CustomerUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerUsecase {
	mock := &CustomerUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
