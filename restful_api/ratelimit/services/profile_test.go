package services

import (
	"context"
	"fmt"
	"testing"
	"time"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	mockRedis "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache/mocks"
	mockDB "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db/mocks"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func Test_services_CreateUserProfile(t *testing.T) {

	type fields struct {
		redis              *mockRedis.MockRedisRepository
		profilesRepository *mockDB.MockProfilesRepository
		otpRepository      *mockDB.MockOtpRepository
	}

	type args struct {
		ctx context.Context
		req models.ProifleRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		fields  fields
		args    args
		senario func(args args, fields fields)
		want    models.ProifleResponse
	}{
		// TODO: Add test cases.
		{
			name: "create user profile success",
			fields: fields{
				redis:              mockRedis.NewMockRedisRepository(ctrl),
				profilesRepository: mockDB.NewMockProfilesRepository(ctrl),
				otpRepository:      mockDB.NewMockOtpRepository(ctrl),
			},
			args: args{
				ctx: context.Background(),
				req: models.ProifleRequest{
					FirstName: "watcharis",
					LastName:  "sukcha",
					Email:     "xxx@test.com",
					Phone:     "0812345678",
				},
			},
			senario: func(args args, fields fields) {

				fields.profilesRepository.EXPECT().CreateUserProfile(args.ctx, gomock.Any()).Return(models.ProfileDB{
					ID:        1,
					UUID:      "1234",
					FirstName: args.req.FirstName,
					LastName:  args.req.LastName,
					Email:     args.req.Email,
					Phone:     args.req.Phone,
					CreatedAt: time.Date(2021, 10, 10, 10, 10, 10, 10, time.UTC),
					UpdatedAt: time.Date(2021, 10, 10, 10, 10, 10, 10, time.UTC),
				}, nil)

				redisProflekey := fmt.Sprintf(models.REDIS_USER_PROFILE, "1234")
				fields.redis.EXPECT().Hset(args.ctx, redisProflekey, gomock.Any()).Return(int64(1), nil)

			},
			want: models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, 200),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.senario(tt.args, tt.fields)

			s := NewServices(tt.fields.redis, tt.fields.profilesRepository, tt.fields.otpRepository)

			got, err := s.CreateUserProfile(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)

		})
	}
}

func Test_services_GetUserProfile(t *testing.T) {
	type fields struct {
		redis              *mockRedis.MockRedisRepository
		profilesRepository *mockDB.MockProfilesRepository
		otpRepository      *mockDB.MockOtpRepository
	}
	type args struct {
		ctx  context.Context
		uuid string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		fields  fields
		args    args
		senario func(args args, fields fields)
		want    models.ProifleResponse
	}{
		// TODO: Add test cases.
		{
			name: "get user profile success",
			fields: fields{
				redis:              mockRedis.NewMockRedisRepository(ctrl),
				profilesRepository: mockDB.NewMockProfilesRepository(ctrl),
				otpRepository:      mockDB.NewMockOtpRepository(ctrl),
			},
			args: args{
				ctx:  context.Background(),
				uuid: "1234",
			},
			senario: func(args args, fields fields) {

				redisProflekey := fmt.Sprintf(models.REDIS_USER_PROFILE, "1234")
				fields.redis.EXPECT().HgetallProfile(args.ctx, redisProflekey).Return(models.ProfileDB{
					ID:        1,
					UUID:      "1234",
					FirstName: "watcharis",
					LastName:  "sukcha",
					Email:     "xxx@test.com",
					Phone:     "0812345678",
				}, nil)
			},
			want: models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, 200),
				Data: &models.ProfileDB{
					ID:        1,
					UUID:      "1234",
					FirstName: "watcharis",
					LastName:  "sukcha",
					Email:     "xxx@test.com",
					Phone:     "0812345678",
				},
			},
		},
		{
			name: "get user profile from db success",
			fields: fields{
				redis:              mockRedis.NewMockRedisRepository(ctrl),
				profilesRepository: mockDB.NewMockProfilesRepository(ctrl),
				otpRepository:      mockDB.NewMockOtpRepository(ctrl),
			},
			args: args{
				ctx:  context.Background(),
				uuid: "1234",
			},
			senario: func(args args, fields fields) {

				redisProflekey := fmt.Sprintf(models.REDIS_USER_PROFILE, "1234")
				fields.redis.EXPECT().HgetallProfile(args.ctx, redisProflekey).Return(models.ProfileDB{}, redis.Nil)

				fields.profilesRepository.EXPECT().GetUserProfile(args.ctx, args.uuid).Return(models.ProfileDB{
					ID:        1,
					UUID:      "1234",
					FirstName: "watcharis",
					LastName:  "sukcha",
					Email:     "xxx@test.com",
					Phone:     "0812345678",
				}, nil)
			},
			want: models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, 200),
				Data: &models.ProfileDB{
					ID:        1,
					UUID:      "1234",
					FirstName: "watcharis",
					LastName:  "sukcha",
					Email:     "xxx@test.com",
					Phone:     "0812345678",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.senario(tt.args, tt.fields)

			s := NewServices(tt.fields.redis, tt.fields.profilesRepository, tt.fields.otpRepository)

			got, err := s.GetUserProfile(tt.args.ctx, tt.args.uuid)

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}
