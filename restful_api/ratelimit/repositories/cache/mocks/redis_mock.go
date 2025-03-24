// Code generated by MockGen. DO NOT EDIT.
// Source: cache/redis.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"
	models "watcharis/go-poc-protocal/restful_api/ratelimit/models"

	gomock "github.com/golang/mock/gomock"
)

// MockRedisRepository is a mock of RedisRepository interface.
type MockRedisRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRedisRepositoryMockRecorder
}

// MockRedisRepositoryMockRecorder is the mock recorder for MockRedisRepository.
type MockRedisRepositoryMockRecorder struct {
	mock *MockRedisRepository
}

// NewMockRedisRepository creates a new mock instance.
func NewMockRedisRepository(ctrl *gomock.Controller) *MockRedisRepository {
	mock := &MockRedisRepository{ctrl: ctrl}
	mock.recorder = &MockRedisRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisRepository) EXPECT() *MockRedisRepositoryMockRecorder {
	return m.recorder
}

// Expire mocks base method.
func (m *MockRedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expire", ctx, key, expiration)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expire indicates an expected call of Expire.
func (mr *MockRedisRepositoryMockRecorder) Expire(ctx, key, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expire", reflect.TypeOf((*MockRedisRepository)(nil).Expire), ctx, key, expiration)
}

// Get mocks base method.
func (m *MockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisRepositoryMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisRepository)(nil).Get), ctx, key)
}

// Hgetall mocks base method.
func (m *MockRedisRepository) Hgetall(ctx context.Context, key string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hgetall", ctx, key)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hgetall indicates an expected call of Hgetall.
func (mr *MockRedisRepositoryMockRecorder) Hgetall(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hgetall", reflect.TypeOf((*MockRedisRepository)(nil).Hgetall), ctx, key)
}

// HgetallProfile mocks base method.
func (m *MockRedisRepository) HgetallProfile(ctx context.Context, key string) (models.ProfileDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HgetallProfile", ctx, key)
	ret0, _ := ret[0].(models.ProfileDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HgetallProfile indicates an expected call of HgetallProfile.
func (mr *MockRedisRepositoryMockRecorder) HgetallProfile(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HgetallProfile", reflect.TypeOf((*MockRedisRepository)(nil).HgetallProfile), ctx, key)
}

// Hset mocks base method.
func (m *MockRedisRepository) Hset(ctx context.Context, key string, values []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hset", ctx, key, values)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hset indicates an expected call of Hset.
func (mr *MockRedisRepositoryMockRecorder) Hset(ctx, key, values interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hset", reflect.TypeOf((*MockRedisRepository)(nil).Hset), ctx, key, values)
}

// Increment mocks base method.
func (m *MockRedisRepository) Increment(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increment", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Increment indicates an expected call of Increment.
func (mr *MockRedisRepositoryMockRecorder) Increment(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockRedisRepository)(nil).Increment), ctx, key)
}

// Set mocks base method.
func (m *MockRedisRepository) Set(ctx context.Context, key, value string, expiration time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiration)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockRedisRepositoryMockRecorder) Set(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisRepository)(nil).Set), ctx, key, value, expiration)
}
