// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	domain "kreditplus/internal/domain"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// LimitRepository is an autogenerated mock type for the LimitRepository type
type LimitRepository struct {
	mock.Mock
}

// CreateLimit provides a mock function with given fields: limit
func (_m *LimitRepository) CreateLimit(limit *domain.Limit) error {
	ret := _m.Called(limit)

	if len(ret) == 0 {
		panic("no return value specified for CreateLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Limit) error); ok {
		r0 = rf(limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteLimit provides a mock function with given fields: id
func (_m *LimitRepository) DeleteLimit(id uint) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllLimits provides a mock function with given fields: limit, offset
func (_m *LimitRepository) GetAllLimits(limit int, offset int) ([]domain.Limit, error) {
	ret := _m.Called(limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetAllLimits")
	}

	var r0 []domain.Limit
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]domain.Limit, error)); ok {
		return rf(limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int, int) []domain.Limit); ok {
		r0 = rf(limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Limit)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLimitByID provides a mock function with given fields: id
func (_m *LimitRepository) GetLimitByID(id uint) (*domain.Limit, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetLimitByID")
	}

	var r0 *domain.Limit
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*domain.Limit, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *domain.Limit); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Limit)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLimitByIDWithTx provides a mock function with given fields: tx, id
func (_m *LimitRepository) GetLimitByIDWithTx(tx *gorm.DB, id uint) (*domain.Limit, error) {
	ret := _m.Called(tx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetLimitByIDWithTx")
	}

	var r0 *domain.Limit
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) (*domain.Limit, error)); ok {
		return rf(tx, id)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) *domain.Limit); ok {
		r0 = rf(tx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Limit)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLimitByNIKandTenorWithTx provides a mock function with given fields: tx, nik, tenor
func (_m *LimitRepository) GetLimitByNIKandTenorWithTx(tx *gorm.DB, nik string, tenor float64) (*domain.Limit, error) {
	ret := _m.Called(tx, nik, tenor)

	if len(ret) == 0 {
		panic("no return value specified for GetLimitByNIKandTenorWithTx")
	}

	var r0 *domain.Limit
	var r1 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, float64) (*domain.Limit, error)); ok {
		return rf(tx, nik, tenor)
	}
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, float64) *domain.Limit); ok {
		r0 = rf(tx, nik, tenor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Limit)
		}
	}

	if rf, ok := ret.Get(1).(func(*gorm.DB, string, float64) error); ok {
		r1 = rf(tx, nik, tenor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLimit provides a mock function with given fields: limit
func (_m *LimitRepository) UpdateLimit(limit *domain.Limit) error {
	ret := _m.Called(limit)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Limit) error); ok {
		r0 = rf(limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLimitWithTx provides a mock function with given fields: tx, limit
func (_m *LimitRepository) UpdateLimitWithTx(tx *gorm.DB, limit *domain.Limit) error {
	ret := _m.Called(tx, limit)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLimitWithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *domain.Limit) error); ok {
		r0 = rf(tx, limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLimitRepository creates a new instance of LimitRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLimitRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *LimitRepository {
	mock := &LimitRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
