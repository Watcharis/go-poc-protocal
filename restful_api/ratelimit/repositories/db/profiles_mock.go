package db

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/stretchr/testify/mock"
)

type ProfilesRepositoryMock struct {
	mock.Mock
}

func (m *ProfilesRepositoryMock) CreateUserProfile(ctx context.Context, data models.ProfileDB) (models.ProfileDB, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(models.ProfileDB), args.Error(1)
}

func (m *ProfilesRepositoryMock) GetUserProfile(ctx context.Context, uuid string) (models.ProfileDB, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(models.ProfileDB), args.Error(1)
}
