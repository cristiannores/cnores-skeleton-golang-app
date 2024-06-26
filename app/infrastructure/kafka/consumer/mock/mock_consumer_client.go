// Code generated by MockGen. DO NOT EDIT.
// Source: consumer_client.go

// Package mock_consumer is a generated GoMock package.
package mock

import (
	context "context"
	"cnores-skeleton-golang-app/app/infrastructure/kafka/consumer"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockConsumerClientInterface is a mock of ConsumerClientInterface interfaces.
type MockConsumerClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerClientInterfaceMockRecorder
}

// MockConsumerClientInterfaceMockRecorder is the mock recorder for MockConsumerClientInterface.
type MockConsumerClientInterfaceMockRecorder struct {
	mock *MockConsumerClientInterface
}

// NewMockConsumerClientInterface creates a new mock instance.
func NewMockConsumerClientInterface(ctrl *gomock.Controller) *MockConsumerClientInterface {
	mock := &MockConsumerClientInterface{ctrl: ctrl}
	mock.recorder = &MockConsumerClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerClientInterface) EXPECT() *MockConsumerClientInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConsumerClientInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConsumerClientInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConsumerClientInterface)(nil).Close))
}

// Consumer mocks base method.
func (m *MockConsumerClientInterface) Consumer(ctx context.Context) <-chan consumer.IncomingMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consumer", ctx)
	ret0, _ := ret[0].(<-chan consumer.IncomingMessage)
	return ret0
}

// Consumer indicates an expected call of Consumer.
func (mr *MockConsumerClientInterfaceMockRecorder) Consumer(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consumer", reflect.TypeOf((*MockConsumerClientInterface)(nil).Consumer), ctx)
}
