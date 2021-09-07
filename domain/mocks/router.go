// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	messaging "github.com/sonalys/letterme/domain/messaging"
	mock "github.com/stretchr/testify/mock"
)

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// AddHandler provides a mock function with given fields: eventType, handler
func (_m *Router) AddHandler(eventType messaging.Event, handler messaging.RouterHandler) {
	_m.Called(eventType, handler)
}

// Communicate provides a mock function with given fields: queue, m, dst
func (_m *Router) Communicate(queue messaging.Queue, m messaging.Message, dst interface{}) error {
	ret := _m.Called(queue, m, dst)

	var r0 error
	if rf, ok := ret.Get(0).(func(messaging.Queue, messaging.Message, interface{}) error); ok {
		r0 = rf(queue, m, dst)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
