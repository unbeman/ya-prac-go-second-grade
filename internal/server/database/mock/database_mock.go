// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/unbeman/ya-prac-go-second-grade/internal/server/database (interfaces: Database)

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	model "github.com/unbeman/ya-prac-go-second-grade/internal/server/model"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockDatabase) CreateUser(arg0 context.Context, arg1 model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockDatabaseMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDatabase)(nil).CreateUser), arg0, arg1)
}

// DeleteUserSecrets mocks base method.
func (m *MockDatabase) DeleteUserSecrets(arg0 context.Context, arg1 model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserSecrets", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserSecrets indicates an expected call of DeleteUserSecrets.
func (mr *MockDatabaseMockRecorder) DeleteUserSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserSecrets", reflect.TypeOf((*MockDatabase)(nil).DeleteUserSecrets), arg0, arg1)
}

// GetUserByID mocks base method.
func (m *MockDatabase) GetUserByID(arg0 context.Context, arg1 uuid.UUID) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockDatabaseMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockDatabase)(nil).GetUserByID), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockDatabase) GetUserByLogin(arg0 context.Context, arg1 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockDatabaseMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockDatabase)(nil).GetUserByLogin), arg0, arg1)
}

// GetUserSecrets mocks base method.
func (m *MockDatabase) GetUserSecrets(arg0 context.Context, arg1 model.User) ([]model.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSecrets", arg0, arg1)
	ret0, _ := ret[0].([]model.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSecrets indicates an expected call of GetUserSecrets.
func (mr *MockDatabaseMockRecorder) GetUserSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSecrets", reflect.TypeOf((*MockDatabase)(nil).GetUserSecrets), arg0, arg1)
}

// SaveUserSecrets mocks base method.
func (m *MockDatabase) SaveUserSecrets(arg0 context.Context, arg1 model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserSecrets", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveUserSecrets indicates an expected call of SaveUserSecrets.
func (mr *MockDatabaseMockRecorder) SaveUserSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserSecrets", reflect.TypeOf((*MockDatabase)(nil).SaveUserSecrets), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockDatabase) UpdateUser(arg0 context.Context, arg1 model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockDatabaseMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockDatabase)(nil).UpdateUser), arg0, arg1)
}
