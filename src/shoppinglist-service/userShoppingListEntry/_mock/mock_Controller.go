// Code generated by mockery v2.39.1. DO NOT EDIT.

package userShoppingListEntry

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockController is an autogenerated mock type for the Controller type
type MockController struct {
	mock.Mock
}

type MockController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockController) EXPECT() *MockController_Expecter {
	return &MockController_Expecter{mock: &_m.Mock}
}

// DeleteEntry provides a mock function with given fields: _a0, _a1
func (_m *MockController) DeleteEntry(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockController_DeleteEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteEntry'
type MockController_DeleteEntry_Call struct {
	*mock.Call
}

// DeleteEntry is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockController_Expecter) DeleteEntry(_a0 interface{}, _a1 interface{}) *MockController_DeleteEntry_Call {
	return &MockController_DeleteEntry_Call{Call: _e.mock.On("DeleteEntry", _a0, _a1)}
}

func (_c *MockController_DeleteEntry_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockController_DeleteEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockController_DeleteEntry_Call) Return() *MockController_DeleteEntry_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockController_DeleteEntry_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockController_DeleteEntry_Call {
	_c.Call.Return(run)
	return _c
}

// GetEntries provides a mock function with given fields: _a0, _a1
func (_m *MockController) GetEntries(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockController_GetEntries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEntries'
type MockController_GetEntries_Call struct {
	*mock.Call
}

// GetEntries is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockController_Expecter) GetEntries(_a0 interface{}, _a1 interface{}) *MockController_GetEntries_Call {
	return &MockController_GetEntries_Call{Call: _e.mock.On("GetEntries", _a0, _a1)}
}

func (_c *MockController_GetEntries_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockController_GetEntries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockController_GetEntries_Call) Return() *MockController_GetEntries_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockController_GetEntries_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockController_GetEntries_Call {
	_c.Call.Return(run)
	return _c
}

// GetEntry provides a mock function with given fields: _a0, _a1
func (_m *MockController) GetEntry(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockController_GetEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEntry'
type MockController_GetEntry_Call struct {
	*mock.Call
}

// GetEntry is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockController_Expecter) GetEntry(_a0 interface{}, _a1 interface{}) *MockController_GetEntry_Call {
	return &MockController_GetEntry_Call{Call: _e.mock.On("GetEntry", _a0, _a1)}
}

func (_c *MockController_GetEntry_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockController_GetEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockController_GetEntry_Call) Return() *MockController_GetEntry_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockController_GetEntry_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockController_GetEntry_Call {
	_c.Call.Return(run)
	return _c
}

// PostEntry provides a mock function with given fields: _a0, _a1
func (_m *MockController) PostEntry(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockController_PostEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PostEntry'
type MockController_PostEntry_Call struct {
	*mock.Call
}

// PostEntry is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockController_Expecter) PostEntry(_a0 interface{}, _a1 interface{}) *MockController_PostEntry_Call {
	return &MockController_PostEntry_Call{Call: _e.mock.On("PostEntry", _a0, _a1)}
}

func (_c *MockController_PostEntry_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockController_PostEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockController_PostEntry_Call) Return() *MockController_PostEntry_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockController_PostEntry_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockController_PostEntry_Call {
	_c.Call.Return(run)
	return _c
}

// PutEntry provides a mock function with given fields: _a0, _a1
func (_m *MockController) PutEntry(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}

// MockController_PutEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutEntry'
type MockController_PutEntry_Call struct {
	*mock.Call
}

// PutEntry is a helper method to define mock.On call
//   - _a0 http.ResponseWriter
//   - _a1 *http.Request
func (_e *MockController_Expecter) PutEntry(_a0 interface{}, _a1 interface{}) *MockController_PutEntry_Call {
	return &MockController_PutEntry_Call{Call: _e.mock.On("PutEntry", _a0, _a1)}
}

func (_c *MockController_PutEntry_Call) Run(run func(_a0 http.ResponseWriter, _a1 *http.Request)) *MockController_PutEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request))
	})
	return _c
}

func (_c *MockController_PutEntry_Call) Return() *MockController_PutEntry_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockController_PutEntry_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request)) *MockController_PutEntry_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockController creates a new instance of MockController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockController {
	mock := &MockController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
