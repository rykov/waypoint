// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// isRunnerJobStreamResponse_Event is an autogenerated mock type for the isRunnerJobStreamResponse_Event type
type isRunnerJobStreamResponse_Event struct {
	mock.Mock
}

// isRunnerJobStreamResponse_Event provides a mock function with given fields:
func (_m *isRunnerJobStreamResponse_Event) isRunnerJobStreamResponse_Event() {
	_m.Called()
}

type mockConstructorTestingTnewIsRunnerJobStreamResponse_Event interface {
	mock.TestingT
	Cleanup(func())
}

// newIsRunnerJobStreamResponse_Event creates a new instance of isRunnerJobStreamResponse_Event. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newIsRunnerJobStreamResponse_Event(t mockConstructorTestingTnewIsRunnerJobStreamResponse_Event) *isRunnerJobStreamResponse_Event {
	mock := &isRunnerJobStreamResponse_Event{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
