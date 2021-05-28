// Code generated by MockGen. DO NOT EDIT.
// Source: ./../storage/types.go

// Package storageMocks is a generated GoMock package.
package storageMocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/lidofinance/dc4bc/storage"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetMessages mocks base method.
func (m *MockStorage) GetMessages(offset uint64) ([]storage.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessages", offset)
	ret0, _ := ret[0].([]storage.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessages indicates an expected call of GetMessages.
func (mr *MockStorageMockRecorder) GetMessages(offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessages", reflect.TypeOf((*MockStorage)(nil).GetMessages), offset)
}

// IgnoreMessages mocks base method.
func (m *MockStorage) IgnoreMessages(messages []string, useOffset bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IgnoreMessages", messages, useOffset)
	ret0, _ := ret[0].(error)
	return ret0
}

// IgnoreMessages indicates an expected call of IgnoreMessages.
func (mr *MockStorageMockRecorder) IgnoreMessages(messages, useOffset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IgnoreMessages", reflect.TypeOf((*MockStorage)(nil).IgnoreMessages), messages, useOffset)
}

// Send mocks base method.
func (m *MockStorage) Send(messages ...storage.Message) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range messages {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockStorageMockRecorder) Send(messages ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockStorage)(nil).Send), messages...)
}

// UnignoreMessages mocks base method.
func (m *MockStorage) UnignoreMessages() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnignoreMessages")
}

// UnignoreMessages indicates an expected call of UnignoreMessages.
func (mr *MockStorageMockRecorder) UnignoreMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnignoreMessages", reflect.TypeOf((*MockStorage)(nil).UnignoreMessages))
}
