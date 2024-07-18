// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sunnyegg/go-so/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "github.com/sunnyegg/go-so/db/sqlc"
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

// CreateAttendanceMember mocks base method.
func (m *MockStore) CreateAttendanceMember(arg0 context.Context, arg1 db.CreateAttendanceMemberParams) (db.AttendanceMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAttendanceMember", arg0, arg1)
	ret0, _ := ret[0].(db.AttendanceMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAttendanceMember indicates an expected call of CreateAttendanceMember.
func (mr *MockStoreMockRecorder) CreateAttendanceMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAttendanceMember", reflect.TypeOf((*MockStore)(nil).CreateAttendanceMember), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateStream mocks base method.
func (m *MockStore) CreateStream(arg0 context.Context, arg1 db.CreateStreamParams) (db.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStream", arg0, arg1)
	ret0, _ := ret[0].(db.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStream indicates an expected call of CreateStream.
func (mr *MockStoreMockRecorder) CreateStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStream", reflect.TypeOf((*MockStore)(nil).CreateStream), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserConfig mocks base method.
func (m *MockStore) CreateUserConfig(arg0 context.Context, arg1 db.CreateUserConfigParams) (db.UserConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserConfig", arg0, arg1)
	ret0, _ := ret[0].(db.UserConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserConfig indicates an expected call of CreateUserConfig.
func (mr *MockStoreMockRecorder) CreateUserConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserConfig", reflect.TypeOf((*MockStore)(nil).CreateUserConfig), arg0, arg1)
}

// DeleteStream mocks base method.
func (m *MockStore) DeleteStream(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStream", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStream indicates an expected call of DeleteStream.
func (mr *MockStoreMockRecorder) DeleteStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStream", reflect.TypeOf((*MockStore)(nil).DeleteStream), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// DeleteUserConfig mocks base method.
func (m *MockStore) DeleteUserConfig(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserConfig indicates an expected call of DeleteUserConfig.
func (mr *MockStoreMockRecorder) DeleteUserConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserConfig", reflect.TypeOf((*MockStore)(nil).DeleteUserConfig), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 db.GetSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetStream mocks base method.
func (m *MockStore) GetStream(arg0 context.Context, arg1 db.GetStreamParams) (db.GetStreamRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStream", arg0, arg1)
	ret0, _ := ret[0].(db.GetStreamRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStream indicates an expected call of GetStream.
func (mr *MockStoreMockRecorder) GetStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStream", reflect.TypeOf((*MockStore)(nil).GetStream), arg0, arg1)
}

// GetStreamAttendanceMembers mocks base method.
func (m *MockStore) GetStreamAttendanceMembers(arg0 context.Context, arg1 db.GetStreamAttendanceMembersParams) ([]db.GetStreamAttendanceMembersRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStreamAttendanceMembers", arg0, arg1)
	ret0, _ := ret[0].([]db.GetStreamAttendanceMembersRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStreamAttendanceMembers indicates an expected call of GetStreamAttendanceMembers.
func (mr *MockStoreMockRecorder) GetStreamAttendanceMembers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStreamAttendanceMembers", reflect.TypeOf((*MockStore)(nil).GetStreamAttendanceMembers), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 int64) (db.GetUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.GetUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByUserID mocks base method.
func (m *MockStore) GetUserByUserID(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUserID", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUserID indicates an expected call of GetUserByUserID.
func (mr *MockStoreMockRecorder) GetUserByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUserID", reflect.TypeOf((*MockStore)(nil).GetUserByUserID), arg0, arg1)
}

// GetUserConfig mocks base method.
func (m *MockStore) GetUserConfig(arg0 context.Context, arg1 db.GetUserConfigParams) (db.UserConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserConfig", arg0, arg1)
	ret0, _ := ret[0].(db.UserConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserConfig indicates an expected call of GetUserConfig.
func (mr *MockStoreMockRecorder) GetUserConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserConfig", reflect.TypeOf((*MockStore)(nil).GetUserConfig), arg0, arg1)
}

// ListStreams mocks base method.
func (m *MockStore) ListStreams(arg0 context.Context, arg1 db.ListStreamsParams) ([]db.ListStreamsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStreams", arg0, arg1)
	ret0, _ := ret[0].([]db.ListStreamsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStreams indicates an expected call of ListStreams.
func (mr *MockStoreMockRecorder) ListStreams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStreams", reflect.TypeOf((*MockStore)(nil).ListStreams), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}

// UpdateUserConfig mocks base method.
func (m *MockStore) UpdateUserConfig(arg0 context.Context, arg1 db.UpdateUserConfigParams) (db.UserConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserConfig", arg0, arg1)
	ret0, _ := ret[0].(db.UserConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserConfig indicates an expected call of UpdateUserConfig.
func (mr *MockStoreMockRecorder) UpdateUserConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserConfig", reflect.TypeOf((*MockStore)(nil).UpdateUserConfig), arg0, arg1)
}
