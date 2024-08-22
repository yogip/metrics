// Code generated by MockGen. DO NOT EDIT.
// Source: metrics/internal/core/service (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_db_store.go -package=mocks metrics/internal/core/service Store
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "metrics/internal/core/model"

	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// BatchUpsertMetrics mocks base method.
func (m *MockStore) BatchUpsertMetrics(arg0 context.Context, arg1 []*model.MetricsV2) ([]*model.MetricsV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchUpsertMetrics", arg0, arg1)
	ret0, _ := ret[0].([]*model.MetricsV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchUpsertMetrics indicates an expected call of BatchUpsertMetrics.
func (mr *MockStoreMockRecorder) BatchUpsertMetrics(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchUpsertMetrics", reflect.TypeOf((*MockStore)(nil).BatchUpsertMetrics), arg0, arg1)
}

// GetCounter mocks base method.
func (m *MockStore) GetCounter(arg0 context.Context, arg1 *model.MetricsV2) (*model.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", arg0, arg1)
	ret0, _ := ret[0].(*model.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockStoreMockRecorder) GetCounter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockStore)(nil).GetCounter), arg0, arg1)
}

// GetGauge mocks base method.
func (m *MockStore) GetGauge(arg0 context.Context, arg1 *model.MetricsV2) (*model.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", arg0, arg1)
	ret0, _ := ret[0].(*model.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockStoreMockRecorder) GetGauge(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockStore)(nil).GetGauge), arg0, arg1)
}

// ListCounter mocks base method.
func (m *MockStore) ListCounter(arg0 context.Context) ([]*model.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCounter", arg0)
	ret0, _ := ret[0].([]*model.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCounter indicates an expected call of ListCounter.
func (mr *MockStoreMockRecorder) ListCounter(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCounter", reflect.TypeOf((*MockStore)(nil).ListCounter), arg0)
}

// ListGauge mocks base method.
func (m *MockStore) ListGauge(arg0 context.Context) ([]*model.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGauge", arg0)
	ret0, _ := ret[0].([]*model.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGauge indicates an expected call of ListGauge.
func (mr *MockStoreMockRecorder) ListGauge(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGauge", reflect.TypeOf((*MockStore)(nil).ListGauge), arg0)
}

// SetCounter mocks base method.
func (m *MockStore) SetCounter(arg0 context.Context, arg1 *model.Counter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCounter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCounter indicates an expected call of SetCounter.
func (mr *MockStoreMockRecorder) SetCounter(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCounter", reflect.TypeOf((*MockStore)(nil).SetCounter), arg0, arg1)
}

// SetGauge mocks base method.
func (m *MockStore) SetGauge(arg0 context.Context, arg1 *model.Gauge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetGauge", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetGauge indicates an expected call of SetGauge.
func (mr *MockStoreMockRecorder) SetGauge(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGauge", reflect.TypeOf((*MockStore)(nil).SetGauge), arg0, arg1)
}
