// Code generated by mockery v2.46.3. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	services "github.com/a-novel/uservice-passkeys/pkg/services"
	mock "github.com/stretchr/testify/mock"
)

// MockGetPasskey is an autogenerated mock type for the GetPasskey type
type MockGetPasskey struct {
	mock.Mock
}

type MockGetPasskey_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGetPasskey) EXPECT() *MockGetPasskey_Expecter {
	return &MockGetPasskey_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockGetPasskey) Exec(ctx context.Context, data *services.GetPasskeyRequest) (*services.GetPasskeyResponse, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *services.GetPasskeyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *services.GetPasskeyRequest) (*services.GetPasskeyResponse, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *services.GetPasskeyRequest) *services.GetPasskeyResponse); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.GetPasskeyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *services.GetPasskeyRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGetPasskey_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockGetPasskey_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data *services.GetPasskeyRequest
func (_e *MockGetPasskey_Expecter) Exec(ctx interface{}, data interface{}) *MockGetPasskey_Exec_Call {
	return &MockGetPasskey_Exec_Call{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockGetPasskey_Exec_Call) Run(run func(ctx context.Context, data *services.GetPasskeyRequest)) *MockGetPasskey_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*services.GetPasskeyRequest))
	})
	return _c
}

func (_c *MockGetPasskey_Exec_Call) Return(_a0 *services.GetPasskeyResponse, _a1 error) *MockGetPasskey_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGetPasskey_Exec_Call) RunAndReturn(run func(context.Context, *services.GetPasskeyRequest) (*services.GetPasskeyResponse, error)) *MockGetPasskey_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGetPasskey creates a new instance of MockGetPasskey. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGetPasskey(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGetPasskey {
	mock := &MockGetPasskey{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}