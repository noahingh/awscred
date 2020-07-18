// Code generated by MockGen. DO NOT EDIT.
// Source: internal/daemon/server/usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	fsnotify "github.com/fsnotify/fsnotify"
	gomock "github.com/golang/mock/gomock"
	core "github.com/hanjunlee/awscred/core"
	reflect "reflect"
)

// MockSessionTokenGenerator is a mock of SessionTokenGenerator interface
type MockSessionTokenGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockSessionTokenGeneratorMockRecorder
}

// MockSessionTokenGeneratorMockRecorder is the mock recorder for MockSessionTokenGenerator
type MockSessionTokenGeneratorMockRecorder struct {
	mock *MockSessionTokenGenerator
}

// NewMockSessionTokenGenerator creates a new mock instance
func NewMockSessionTokenGenerator(ctrl *gomock.Controller) *MockSessionTokenGenerator {
	mock := &MockSessionTokenGenerator{ctrl: ctrl}
	mock.recorder = &MockSessionTokenGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessionTokenGenerator) EXPECT() *MockSessionTokenGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockSessionTokenGenerator) Generate(arg0 core.Cred, arg1 core.Config, arg2 string) (core.SessionToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", arg0, arg1, arg2)
	ret0, _ := ret[0].(core.SessionToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (mr *MockSessionTokenGeneratorMockRecorder) Generate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSessionTokenGenerator)(nil).Generate), arg0, arg1, arg2)
}

// MockFileWatcher is a mock of FileWatcher interface
type MockFileWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockFileWatcherMockRecorder
}

// MockFileWatcherMockRecorder is the mock recorder for MockFileWatcher
type MockFileWatcherMockRecorder struct {
	mock *MockFileWatcher
}

// NewMockFileWatcher creates a new mock instance
func NewMockFileWatcher(ctrl *gomock.Controller) *MockFileWatcher {
	mock := &MockFileWatcher{ctrl: ctrl}
	mock.recorder = &MockFileWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileWatcher) EXPECT() *MockFileWatcherMockRecorder {
	return m.recorder
}

// Watch mocks base method
func (m *MockFileWatcher) Watch(arg0 context.Context, arg1 chan<- fsnotify.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0, arg1)
}

// Watch indicates an expected call of Watch
func (mr *MockFileWatcherMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockFileWatcher)(nil).Watch), arg0, arg1)
}

// MockCredFileHandler is a mock of CredFileHandler interface
type MockCredFileHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCredFileHandlerMockRecorder
}

// MockCredFileHandlerMockRecorder is the mock recorder for MockCredFileHandler
type MockCredFileHandlerMockRecorder struct {
	mock *MockCredFileHandler
}

// NewMockCredFileHandler creates a new mock instance
func NewMockCredFileHandler(ctrl *gomock.Controller) *MockCredFileHandler {
	mock := &MockCredFileHandler{ctrl: ctrl}
	mock.recorder = &MockCredFileHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCredFileHandler) EXPECT() *MockCredFileHandlerMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockCredFileHandler) Read() (map[string]core.Cred, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(map[string]core.Cred)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockCredFileHandlerMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockCredFileHandler)(nil).Read))
}

// Write mocks base method
func (m *MockCredFileHandler) Write(arg0 map[string]core.Cred) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockCredFileHandlerMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCredFileHandler)(nil).Write), arg0)
}

// Remove mocks base method
func (m *MockCredFileHandler) Remove() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove")
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockCredFileHandlerMockRecorder) Remove() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockCredFileHandler)(nil).Remove))
}

// MockConfigFileHandler is a mock of ConfigFileHandler interface
type MockConfigFileHandler struct {
	ctrl     *gomock.Controller
	recorder *MockConfigFileHandlerMockRecorder
}

// MockConfigFileHandlerMockRecorder is the mock recorder for MockConfigFileHandler
type MockConfigFileHandlerMockRecorder struct {
	mock *MockConfigFileHandler
}

// NewMockConfigFileHandler creates a new mock instance
func NewMockConfigFileHandler(ctrl *gomock.Controller) *MockConfigFileHandler {
	mock := &MockConfigFileHandler{ctrl: ctrl}
	mock.recorder = &MockConfigFileHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigFileHandler) EXPECT() *MockConfigFileHandlerMockRecorder {
	return m.recorder
}

// Read mocks base method
func (m *MockConfigFileHandler) Read() (map[string]core.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(map[string]core.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockConfigFileHandlerMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConfigFileHandler)(nil).Read))
}

// Write mocks base method
func (m *MockConfigFileHandler) Write(arg0 map[string]core.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockConfigFileHandlerMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockConfigFileHandler)(nil).Write), arg0)
}
