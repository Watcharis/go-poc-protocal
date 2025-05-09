package services

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache"
	mockRedis "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache/mocks"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db"
	mockDB "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db/mocks"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func Test_services_VerifyOtpRatelimit(t *testing.T) {

	type mockRepository struct {
		mockRedis              *cache.RedisRepositoryMock
		mockProfilesRepository *db.ProfilesRepositoryMock
		mockOtpRepository      *db.OtpRepositoryMock
	}

	var initMockRepository = func() mockRepository {
		return mockRepository{
			mockRedis:              &cache.RedisRepositoryMock{},
			mockProfilesRepository: &db.ProfilesRepositoryMock{},
			mockOtpRepository:      &db.OtpRepositoryMock{},
		}
	}

	type args struct {
		ctx context.Context
		req models.VerifyOtpRatelimitRequest
	}

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	tests := []struct {
		name    string
		args    args
		senario func(args args) mockRepository
		want    models.VerifyOtpRatelimitResponse
	}{
		// TODO: Add test cases.
		{
			name: "get otp ratelimit success",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository()

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.On("Get", args.ctx, redisKeyRatelimitOTP).Return("1", nil)

				redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, args.req.Uuid)
				mockRepositories.mockRedis.On("Get", args.ctx, redisKeyOTP).Return("200139", nil)

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				Data: &models.VerifyOtpRatelimitDataResponse{
					Otp: "200139",
				},
				Error: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepositories := tt.senario(tt.args)

			s := NewServices(mockRepositories.mockRedis, mockRepositories.mockProfilesRepository, mockRepositories.mockOtpRepository)

			got, err := s.VerifyOtpRatelimit(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func Test_services_VerifyOtpRatelimit_use_gomock_gen(t *testing.T) {

	type args struct {
		ctx context.Context
		req models.VerifyOtpRatelimitRequest
	}

	type mockRepository struct {
		mockRedis              *mockRedis.MockRedisRepository
		mockProfilesRepository *mockDB.MockProfilesRepository
		mockOtpRepository      *mockDB.MockOtpRepository
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var initMockRepository = func(ctrl *gomock.Controller) mockRepository {
		return mockRepository{
			mockRedis:              mockRedis.NewMockRedisRepository(ctrl),
			mockProfilesRepository: mockDB.NewMockProfilesRepository(ctrl),
			mockOtpRepository:      mockDB.NewMockOtpRepository(ctrl),
		}
	}
	tests := []struct {
		name    string
		args    args
		senario func(args args) mockRepository
		want    models.VerifyOtpRatelimitResponse
	}{
		// TODO: Add test cases.
		{
			name: "get otp ratelimit success",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository(ctrl)

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyRatelimitOTP).Return("1", nil)

				redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyOTP).Return("200139", nil)

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				Data: &models.VerifyOtpRatelimitDataResponse{
					Otp: "200139",
				},
				Error: nil,
			},
		},
		{
			name: "get count otp ratelimit failed",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "99999",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository(ctrl)

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyRatelimitOTP).Return("xxx", fmt.Errorf("error get countotp ratelimit failed"))

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
				Error: &response.ErrorResponse{
					ErrorMessage: "error get countotp ratelimit failed",
				},
			},
		},
		{
			name: "get otp from redis Nil and otp not match",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "99999",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository(ctrl)

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyRatelimitOTP).Return("1", nil)

				redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyOTP).Return("", redis.Nil)

				mockRepositories.mockOtpRepository.EXPECT().GetOtp(args.ctx, args.req.Uuid, args.req.Otp).Return(models.OtpDB{}, fmt.Errorf("not found otp in db"))

				mockRepositories.mockRedis.EXPECT().Increment(args.ctx, redisKeyRatelimitOTP).Return(int64(2), nil)

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusNotFound),
				Error: &response.ErrorResponse{
					ErrorMessage: "otp not match",
				},
			},
		},
		{
			name: "get otp from redis Nil and otp match in db",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository(ctrl)

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyRatelimitOTP).Return("1", nil)

				redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyOTP).Return("", redis.Nil)

				mockRepositories.mockOtpRepository.EXPECT().GetOtp(args.ctx, args.req.Uuid, args.req.Otp).Return(models.OtpDB{
					ID:   1,
					UUID: "1234",
					Otp:  "200139",
				}, nil)

				mockRepositories.mockRedis.EXPECT().Set(args.ctx, redisKeyOTP, args.req.Otp, time.Duration(models.OTP_EXPIRE)).Return("200139", nil)

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				Data: &models.VerifyOtpRatelimitDataResponse{
					Otp: "200139",
				},
			},
		},
		{
			name: "otp not match from request",
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "999999",
				},
			},
			senario: func(args args) mockRepository {

				mockRepositories := initMockRepository(ctrl)

				redisKeyRatelimitOTP := fmt.Sprintf(models.REDIS_RATELIMIT_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyRatelimitOTP).Return("1", nil)

				redisKeyOTP := fmt.Sprintf(models.REDIS_OTP, args.req.Uuid)
				mockRepositories.mockRedis.EXPECT().Get(args.ctx, redisKeyOTP).Return("200139", nil)

				mockRepositories.mockRedis.EXPECT().Increment(args.ctx, redisKeyRatelimitOTP).Return(int64(3), nil)

				mockRepositories.mockRedis.EXPECT().Expire(args.ctx, redisKeyRatelimitOTP, time.Duration(models.RATELIMIT_OTP_EXPIRE)).Return(true, nil)

				return mockRepositories
			},
			want: models.VerifyOtpRatelimitResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusTooManyRequests),
				Error: &response.ErrorResponse{
					ErrorMessage: "otp ratelimit exceed",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepositories := tt.senario(tt.args)

			s := NewServices(mockRepositories.mockRedis, mockRepositories.mockProfilesRepository, mockRepositories.mockOtpRepository)

			got, err := s.VerifyOtpRatelimit(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}
