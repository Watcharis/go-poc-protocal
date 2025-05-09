// Code generated by MockGen. DO NOT EDIT.
// Source: db/profiles.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	models "watcharis/go-poc-protocal/restful_api/ratelimit/models"

	gomock "github.com/golang/mock/gomock"
)

// MockProfilesRepository is a mock of ProfilesRepository interface.
type MockProfilesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProfilesRepositoryMockRecorder
}

// MockProfilesRepositoryMockRecorder is the mock recorder for MockProfilesRepository.
type MockProfilesRepositoryMockRecorder struct {
	mock *MockProfilesRepository
}

// NewMockProfilesRepository creates a new mock instance.
func NewMockProfilesRepository(ctrl *gomock.Controller) *MockProfilesRepository {
	mock := &MockProfilesRepository{ctrl: ctrl}
	mock.recorder = &MockProfilesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfilesRepository) EXPECT() *MockProfilesRepositoryMockRecorder {
	return m.recorder
}

// CreateUserProfile mocks base method.
func (m *MockProfilesRepository) CreateUserProfile(ctx context.Context, data models.ProfileDB) (models.ProfileDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserProfile", ctx, data)
	ret0, _ := ret[0].(models.ProfileDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserProfile indicates an expected call of CreateUserProfile.
func (mr *MockProfilesRepositoryMockRecorder) CreateUserProfile(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserProfile", reflect.TypeOf((*MockProfilesRepository)(nil).CreateUserProfile), ctx, data)
}

// GetUserProfile mocks base method.
func (m *MockProfilesRepository) GetUserProfile(ctx context.Context, uuid string) (models.ProfileDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, uuid)
	ret0, _ := ret[0].(models.ProfileDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockProfilesRepositoryMockRecorder) GetUserProfile(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockProfilesRepository)(nil).GetUserProfile), ctx, uuid)
}
