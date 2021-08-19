// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/sonalys/letterme/domain"

	mock "github.com/stretchr/testify/mock"
)

// Collection is an autogenerated mock type for the Collection type
type Collection struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, documents
func (_m *Collection) Create(ctx context.Context, documents ...interface{}) ([]domain.DatabaseID, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, documents...)
	ret := _m.Called(_ca...)

	var r0 []domain.DatabaseID
	if rf, ok := ret.Get(0).(func(context.Context, ...interface{}) []domain.DatabaseID); ok {
		r0 = rf(ctx, documents...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.DatabaseID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...interface{}) error); ok {
		r1 = rf(ctx, documents...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, filter
func (_m *Collection) Delete(ctx context.Context, filter interface{}) (int64, error) {
	ret := _m.Called(ctx, filter)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) int64); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// First provides a mock function with given fields: ctx, filter, dst
func (_m *Collection) First(ctx context.Context, filter interface{}, dst interface{}) error {
	ret := _m.Called(ctx, filter, dst)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}) error); ok {
		r0 = rf(ctx, filter, dst)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, filter, dst
func (_m *Collection) List(ctx context.Context, filter interface{}, dst interface{}) error {
	ret := _m.Called(ctx, filter, dst)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}) error); ok {
		r0 = rf(ctx, filter, dst)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, filter, update
func (_m *Collection) Update(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	ret := _m.Called(ctx, filter, update)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}) int64); ok {
		r0 = rf(ctx, filter, update)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}) error); ok {
		r1 = rf(ctx, filter, update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
