// Code generated by MockGen. DO NOT EDIT.
// Source: auth/token_generator.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=_mocks/token_generator.go -source=auth/token_generator.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTokenGenerator is a mock of TokenGenerator interface.
type MockTokenGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockTokenGeneratorMockRecorder
}

// MockTokenGeneratorMockRecorder is the mock recorder for MockTokenGenerator.
type MockTokenGeneratorMockRecorder struct {
	mock *MockTokenGenerator
}

// NewMockTokenGenerator creates a new mock instance.
func NewMockTokenGenerator(ctrl *gomock.Controller) *MockTokenGenerator {
	mock := &MockTokenGenerator{ctrl: ctrl}
	mock.recorder = &MockTokenGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenGenerator) EXPECT() *MockTokenGeneratorMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockTokenGenerator) GenerateToken(claims map[string]any) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", claims)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockTokenGeneratorMockRecorder) GenerateToken(claims any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockTokenGenerator)(nil).GenerateToken), claims)
}
