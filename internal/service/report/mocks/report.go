// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/report/report.go

// Package mock_report is a generated GoMock package.
package mock_report

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// MakeReport mocks base method.
func (m *MockService) MakeReport(ctx context.Context, userId int64, timeSince time.Time, timeRangePrefix string) (*model.MessageOut, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeReport", ctx, userId, timeSince, timeRangePrefix)
	ret0, _ := ret[0].(*model.MessageOut)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeReport indicates an expected call of MakeReport.
func (mr *MockServiceMockRecorder) MakeReport(ctx, userId, timeSince, timeRangePrefix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeReport", reflect.TypeOf((*MockService)(nil).MakeReport), ctx, userId, timeSince, timeRangePrefix)
}
