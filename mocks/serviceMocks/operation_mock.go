// Code generated by MockGen. DO NOT EDIT.
// Source: ./../client/services/operation/operation.go

// Package serviceMocks is a generated GoMock package.
package serviceMocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/lidofinance/dc4bc/client/api/dto"
	types "github.com/lidofinance/dc4bc/client/types"
)

// MockOperationService is a mock of OperationService interface.
type MockOperationService struct {
	ctrl     *gomock.Controller
	recorder *MockOperationServiceMockRecorder
}

// MockOperationServiceMockRecorder is the mock recorder for MockOperationService.
type MockOperationServiceMockRecorder struct {
	mock *MockOperationService
}

// NewMockOperationService creates a new mock instance.
func NewMockOperationService(ctrl *gomock.Controller) *MockOperationService {
	mock := &MockOperationService{ctrl: ctrl}
	mock.recorder = &MockOperationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOperationService) EXPECT() *MockOperationServiceMockRecorder {
	return m.recorder
}

// GetOperation mocks base method.
func (m *MockOperationService) GetOperation(dto *dto.OperationIdDTO) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperation", dto)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperation indicates an expected call of GetOperation.
func (mr *MockOperationServiceMockRecorder) GetOperation(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperation", reflect.TypeOf((*MockOperationService)(nil).GetOperation), dto)
}

// GetOperations mocks base method.
func (m *MockOperationService) GetOperations() (map[string]*types.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperations")
	ret0, _ := ret[0].(map[string]*types.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperations indicates an expected call of GetOperations.
func (mr *MockOperationServiceMockRecorder) GetOperations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperations", reflect.TypeOf((*MockOperationService)(nil).GetOperations))
}
