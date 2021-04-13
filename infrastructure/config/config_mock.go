// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package config is a generated GoMock package.
package config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// GetString mocks base method.
func (m *MockConfig) GetString(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString.
func (mr *MockConfigMockRecorder) GetString(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockConfig)(nil).GetString), key)
}

// IsSet mocks base method.
func (m *MockConfig) IsSet(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSet", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSet indicates an expected call of IsSet.
func (mr *MockConfigMockRecorder) IsSet(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSet", reflect.TypeOf((*MockConfig)(nil).IsSet), key)
}

// Save mocks base method.
func (m *MockConfig) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockConfigMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockConfig)(nil).Save))
}

// Set mocks base method.
func (m *MockConfig) Set(key string, value interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value)
}

// Set indicates an expected call of Set.
func (mr *MockConfigMockRecorder) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockConfig)(nil).Set), key, value)
}

// WriteConfig mocks base method.
func (m *MockConfig) WriteConfig() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteConfig")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteConfig indicates an expected call of WriteConfig.
func (mr *MockConfigMockRecorder) WriteConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteConfig", reflect.TypeOf((*MockConfig)(nil).WriteConfig))
}

// MockLoader is a mock of Loader interface.
type MockLoader struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderMockRecorder
}

// MockLoaderMockRecorder is the mock recorder for MockLoader.
type MockLoaderMockRecorder struct {
	mock *MockLoader
}

// NewMockLoader creates a new mock instance.
func NewMockLoader(ctrl *gomock.Controller) *MockLoader {
	mock := &MockLoader{ctrl: ctrl}
	mock.recorder = &MockLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoader) EXPECT() *MockLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockLoader) Load(name string) (Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", name)
	ret0, _ := ret[0].(Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockLoaderMockRecorder) Load(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockLoader)(nil).Load), name)
}
