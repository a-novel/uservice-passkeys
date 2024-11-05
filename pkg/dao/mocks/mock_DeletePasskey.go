// Code generated by mockery v2.46.3. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/a-novel/uservice-passkeys/pkg/dao"
	entities "github.com/a-novel/uservice-passkeys/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockDeletePasskey is an autogenerated mock type for the DeletePasskey type
type MockDeletePasskey struct {
	mock.Mock
}

type MockDeletePasskey_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeletePasskey) EXPECT() *MockDeletePasskey_Expecter {
	return &MockDeletePasskey_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, request
func (_m *MockDeletePasskey) Exec(ctx context.Context, request *dao.DeletePasskeyRequest) (*entities.Passkey, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *entities.Passkey
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dao.DeletePasskeyRequest) (*entities.Passkey, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dao.DeletePasskeyRequest) *entities.Passkey); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Passkey)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dao.DeletePasskeyRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeletePasskey_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockDeletePasskey_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - request *dao.DeletePasskeyRequest
func (_e *MockDeletePasskey_Expecter) Exec(ctx interface{}, request interface{}) *MockDeletePasskey_Exec_Call {
	return &MockDeletePasskey_Exec_Call{Call: _e.mock.On("Exec", ctx, request)}
}

func (_c *MockDeletePasskey_Exec_Call) Run(run func(ctx context.Context, request *dao.DeletePasskeyRequest)) *MockDeletePasskey_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dao.DeletePasskeyRequest))
	})
	return _c
}

func (_c *MockDeletePasskey_Exec_Call) Return(_a0 *entities.Passkey, _a1 error) *MockDeletePasskey_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeletePasskey_Exec_Call) RunAndReturn(run func(context.Context, *dao.DeletePasskeyRequest) (*entities.Passkey, error)) *MockDeletePasskey_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeletePasskey creates a new instance of MockDeletePasskey. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeletePasskey(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeletePasskey {
	mock := &MockDeletePasskey{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
