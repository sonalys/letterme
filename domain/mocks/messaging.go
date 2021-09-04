// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/sonalys/letterme/domain/models"
)

// Messaging is an autogenerated mock type for the Messaging type
type Messaging struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Messaging) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: ctx, queue, handler
func (_m *Messaging) Consume(ctx context.Context, queue string, handler models.DeliveryHandler) error {
	ret := _m.Called(ctx, queue, handler)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.DeliveryHandler) error); ok {
		r0 = rf(ctx, queue, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateQueue provides a mock function with given fields: name
func (_m *Messaging) CreateQueue(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Publish provides a mock function with given fields: queue, m
func (_m *Messaging) Publish(queue string, m models.Message) error {
	ret := _m.Called(queue, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, models.Message) error); ok {
		r0 = rf(queue, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
