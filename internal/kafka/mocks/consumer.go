// Code generated by MockGen. DO NOT EDIT.
// Source: internal/kafka/consumer.go

// Package mock_kafka is a generated GoMock package.
package mock_kafka

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// HandleMessage mocks base method.
func (m *MockHandler) HandleMessage(ctx context.Context, msg model.ReportMsg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleMessage", ctx, msg)
}

// HandleMessage indicates an expected call of HandleMessage.
func (mr *MockHandlerMockRecorder) HandleMessage(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMessage", reflect.TypeOf((*MockHandler)(nil).HandleMessage), ctx, msg)
}
