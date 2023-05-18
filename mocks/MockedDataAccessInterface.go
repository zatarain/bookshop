// Code generated by mockery v2.27.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// MockedDataAccessInterface is an autogenerated mock type for the MockedDataAccessInterface type
type MockedDataAccessInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *MockedDataAccessInterface) Create(_a0 interface{}) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(interface{}) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

type mockConstructorTestingTNewMockedDataAccessInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockedDataAccessInterface creates a new instance of MockedDataAccessInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockedDataAccessInterface(t mockConstructorTestingTNewMockedDataAccessInterface) *MockedDataAccessInterface {
	mock := &MockedDataAccessInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
