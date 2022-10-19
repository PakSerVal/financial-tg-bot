// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/spend/repository.go

// Package mock_spend is a generated GoMock package.
package mock_spend

import (
	context "context"
	sql "database/sql"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetByTimeSince mocks base method.
func (m *MockRepository) GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTimeSince", ctx, userId, timeSince)
	ret0, _ := ret[0].([]model.Spend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTimeSince indicates an expected call of GetByTimeSince.
func (mr *MockRepositoryMockRecorder) GetByTimeSince(ctx, userId, timeSince interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTimeSince", reflect.TypeOf((*MockRepository)(nil).GetByTimeSince), ctx, userId, timeSince)
}

// GetByTimeSinceTx mocks base method.
func (m *MockRepository) GetByTimeSinceTx(tx *sql.Tx, ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTimeSinceTx", tx, ctx, userId, timeSince)
	ret0, _ := ret[0].([]model.Spend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTimeSinceTx indicates an expected call of GetByTimeSinceTx.
func (mr *MockRepositoryMockRecorder) GetByTimeSinceTx(tx, ctx, userId, timeSince interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTimeSinceTx", reflect.TypeOf((*MockRepository)(nil).GetByTimeSinceTx), tx, ctx, userId, timeSince)
}

// SaveTx mocks base method.
func (m *MockRepository) SaveTx(tx *sql.Tx, ctx context.Context, sum int64, category string, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTx", tx, ctx, sum, category, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveTx indicates an expected call of SaveTx.
func (mr *MockRepositoryMockRecorder) SaveTx(tx, ctx, sum, category, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTx", reflect.TypeOf((*MockRepository)(nil).SaveTx), tx, ctx, sum, category, userId)
}