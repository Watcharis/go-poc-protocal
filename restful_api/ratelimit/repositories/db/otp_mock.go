package db

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"github.com/stretchr/testify/mock"
)

type OtpRepositoryMock struct {
	mock.Mock
}

// mockgen -source=repositories/respositories.go -destination=db/mocks/otp_mock.go -package=mocks -mock_names OtpRepository=MockOtpRepository

func (m *OtpRepositoryMock) CreateOtp(ctx context.Context, data models.OtpDB) (models.OtpDB, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(models.OtpDB), args.Error(1)
}

func (m *OtpRepositoryMock) GetOtp(ctx context.Context, uuid string, otp string) (models.OtpDB, error) {
	args := m.Called(ctx, uuid, otp)
	return args.Get(0).(models.OtpDB), args.Error(1)
}
