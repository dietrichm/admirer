// Code generated by MockGen. DO NOT EDIT.
// Source: services.go

// Package domain is a generated GoMock package.
package domain

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// Authenticate mocks base method.
func (m *MockService) Authenticate(code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockServiceMockRecorder) Authenticate(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockService)(nil).Authenticate), code)
}

// Authenticated mocks base method.
func (m *MockService) Authenticated() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticated")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Authenticated indicates an expected call of Authenticated.
func (mr *MockServiceMockRecorder) Authenticated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticated", reflect.TypeOf((*MockService)(nil).Authenticated))
}

// Close mocks base method.
func (m *MockService) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockServiceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockService)(nil).Close))
}

// CreateAuthURL mocks base method.
func (m *MockService) CreateAuthURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAuthURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateAuthURL indicates an expected call of CreateAuthURL.
func (mr *MockServiceMockRecorder) CreateAuthURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAuthURL", reflect.TypeOf((*MockService)(nil).CreateAuthURL))
}

// GetLovedTracks mocks base method.
func (m *MockService) GetLovedTracks() ([]Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLovedTracks")
	ret0, _ := ret[0].([]Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLovedTracks indicates an expected call of GetLovedTracks.
func (mr *MockServiceMockRecorder) GetLovedTracks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLovedTracks", reflect.TypeOf((*MockService)(nil).GetLovedTracks))
}

// GetUsername mocks base method.
func (m *MockService) GetUsername() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsername")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsername indicates an expected call of GetUsername.
func (mr *MockServiceMockRecorder) GetUsername() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsername", reflect.TypeOf((*MockService)(nil).GetUsername))
}

// Name mocks base method.
func (m *MockService) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockServiceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockService)(nil).Name))
}

// MockServiceLoader is a mock of ServiceLoader interface.
type MockServiceLoader struct {
	ctrl     *gomock.Controller
	recorder *MockServiceLoaderMockRecorder
}

// MockServiceLoaderMockRecorder is the mock recorder for MockServiceLoader.
type MockServiceLoaderMockRecorder struct {
	mock *MockServiceLoader
}

// NewMockServiceLoader creates a new mock instance.
func NewMockServiceLoader(ctrl *gomock.Controller) *MockServiceLoader {
	mock := &MockServiceLoader{ctrl: ctrl}
	mock.recorder = &MockServiceLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceLoader) EXPECT() *MockServiceLoaderMockRecorder {
	return m.recorder
}

// ForName mocks base method.
func (m *MockServiceLoader) ForName(serviceName string) (Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForName", serviceName)
	ret0, _ := ret[0].(Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForName indicates an expected call of ForName.
func (mr *MockServiceLoaderMockRecorder) ForName(serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForName", reflect.TypeOf((*MockServiceLoader)(nil).ForName), serviceName)
}

// Names mocks base method.
func (m *MockServiceLoader) Names() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Names")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Names indicates an expected call of Names.
func (mr *MockServiceLoaderMockRecorder) Names() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Names", reflect.TypeOf((*MockServiceLoader)(nil).Names))
}
