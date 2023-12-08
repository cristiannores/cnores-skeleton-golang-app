// Code generated by MockGen. DO NOT EDIT.
// Source: notification_service.go

// Package mock_notification_service is a generated GoMock package.
package mock_notification_service

import (
	context "context"
	reflect "reflect"

	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	gomock "github.com/golang/mock/gomock"
)

// MockNotificationServiceInterface is a mock of NotificationServiceInterface interface.
type MockNotificationServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceInterfaceMockRecorder
}

// MockNotificationServiceInterfaceMockRecorder is the mock recorder for MockNotificationServiceInterface.
type MockNotificationServiceInterfaceMockRecorder struct {
	mock *MockNotificationServiceInterface
}

// NewMockNotificationServiceInterface creates a new mock instance.
func NewMockNotificationServiceInterface(ctrl *gomock.Controller) *MockNotificationServiceInterface {
	mock := &MockNotificationServiceInterface{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceInterface) EXPECT() *MockNotificationServiceInterfaceMockRecorder {
	return m.recorder
}

// Notify mocks base method.
func (m *MockNotificationServiceInterface) Notify(ctx context.Context, data *notification_service_request.SlackBody) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Notify", ctx, data)
}

// Notify indicates an expected call of Notify.
func (mr *MockNotificationServiceInterfaceMockRecorder) Notify(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockNotificationServiceInterface)(nil).Notify), ctx, data)
}