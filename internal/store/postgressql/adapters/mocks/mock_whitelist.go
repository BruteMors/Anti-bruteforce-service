// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/service/whitelist.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	entity "Anti-bruteforce-service/internal/domain/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWhiteListStore is a mock of WhiteListStore interface.
type MockWhiteListStore struct {
	ctrl     *gomock.Controller
	recorder *MockWhiteListStoreMockRecorder
}

// MockWhiteListStoreMockRecorder is the mock recorder for MockWhiteListStore.
type MockWhiteListStoreMockRecorder struct {
	mock *MockWhiteListStore
}

// NewMockWhiteListStore creates a new mock instance.
func NewMockWhiteListStore(ctrl *gomock.Controller) *MockWhiteListStore {
	mock := &MockWhiteListStore{ctrl: ctrl}
	mock.recorder = &MockWhiteListStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWhiteListStore) EXPECT() *MockWhiteListStoreMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockWhiteListStore) Add(prefix, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", prefix, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockWhiteListStoreMockRecorder) Add(prefix, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockWhiteListStore)(nil).Add), prefix, mask)
}

// Get mocks base method.
func (m *MockWhiteListStore) Get() ([]entity.IpNetwork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].([]entity.IpNetwork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockWhiteListStoreMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockWhiteListStore)(nil).Get))
}

// Remove mocks base method.
func (m *MockWhiteListStore) Remove(prefix, mask string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", prefix, mask)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockWhiteListStoreMockRecorder) Remove(prefix, mask interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockWhiteListStore)(nil).Remove), prefix, mask)
}
