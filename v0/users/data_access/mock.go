// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package userdataaccess is a generated GoMock package.
package userdataaccess

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	usermodel "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
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

// CreateUserToDynamoDb mocks base method.
func (m *MockRepository) CreateUserToDynamoDb(ctx context.Context, user *usermodel.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserToDynamoDb", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUserToDynamoDb indicates an expected call of CreateUserToDynamoDb.
func (mr *MockRepositoryMockRecorder) CreateUserToDynamoDb(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserToDynamoDb", reflect.TypeOf((*MockRepository)(nil).CreateUserToDynamoDb), ctx, user)
}

// GetUsersFromDynamoDb mocks base method.
func (m *MockRepository) GetUsersFromDynamoDb(ctx context.Context, limit int32, nextKey string) ([]*usermodel.User, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersFromDynamoDb", ctx, limit, nextKey)
	ret0, _ := ret[0].([]*usermodel.User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUsersFromDynamoDb indicates an expected call of GetUsersFromDynamoDb.
func (mr *MockRepositoryMockRecorder) GetUsersFromDynamoDb(ctx, limit, nextKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersFromDynamoDb", reflect.TypeOf((*MockRepository)(nil).GetUsersFromDynamoDb), ctx, limit, nextKey)
}
