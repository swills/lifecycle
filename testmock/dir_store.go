// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpacks/lifecycle (interfaces: DirStore)

// Package testmock is a generated GoMock package.
package testmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	buildpack "github.com/buildpacks/lifecycle/buildpack"
)

// MockDirStore is a mock of DirStore interface.
type MockDirStore struct {
	ctrl     *gomock.Controller
	recorder *MockDirStoreMockRecorder
}

// MockDirStoreMockRecorder is the mock recorder for MockDirStore.
type MockDirStoreMockRecorder struct {
	mock *MockDirStore
}

// NewMockDirStore creates a new mock instance.
func NewMockDirStore(ctrl *gomock.Controller) *MockDirStore {
	mock := &MockDirStore{ctrl: ctrl}
	mock.recorder = &MockDirStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDirStore) EXPECT() *MockDirStoreMockRecorder {
	return m.recorder
}

// Lookup mocks base method.
func (m *MockDirStore) Lookup(arg0, arg1, arg2 string) (buildpack.BuildModule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lookup", arg0, arg1, arg2)
	ret0, _ := ret[0].(buildpack.BuildModule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup.
func (mr *MockDirStoreMockRecorder) Lookup(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lookup", reflect.TypeOf((*MockDirStore)(nil).Lookup), arg0, arg1, arg2)
}
