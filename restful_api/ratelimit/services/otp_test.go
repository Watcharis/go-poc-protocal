package services

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	mockRedis "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache/mocks"
	mockDB "watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_services_CreateOtp(t *testing.T) {
	type fields struct {
		redis              *mockRedis.MockRedisRepository
		profilesRepository *mockDB.MockProfilesRepository
		otpRepository      *mockDB.MockOtpRepository
	}
	type args struct {
		ctx context.Context
		req models.OtpRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TimeNow = func() time.Time {
		return time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		senario func(args args, fields fields)
		want    models.OtpResponse
	}{
		// TODO: Add test cases.
		{
			name: "create otp success",
			fields: fields{
				redis:              mockRedis.NewMockRedisRepository(ctrl),
				profilesRepository: mockDB.NewMockProfilesRepository(ctrl),
				otpRepository:      mockDB.NewMockOtpRepository(ctrl),
			},
			args: args{
				ctx: context.Background(),
				req: models.OtpRequest{
					UUID: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args, fields fields) {

				fields.otpRepository.EXPECT().CreateOtp(args.ctx, models.OtpDB{
					UUID:      args.req.UUID,
					Otp:       args.req.Otp,
					CreatedAt: TimeNow(),
				}).Return(models.OtpDB{
					ID:        1,
					UUID:      args.req.UUID,
					Otp:       args.req.Otp,
					CreatedAt: TimeNow(),
					UpdatedAt: TimeNow(),
				}, nil)

				redisKeyOtp := fmt.Sprintf(models.REDIS_OTP, args.req.UUID)
				fields.redis.EXPECT().Set(args.ctx, redisKeyOtp, args.req.Otp, time.Duration(models.OTP_EXPIRE)).Return("1", nil)
			},
			want: models.OtpResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				Error:          nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.senario(tt.args, tt.fields)

			s := NewServices(tt.fields.redis, tt.fields.profilesRepository, tt.fields.otpRepository)

			got, err := s.CreateOtp(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)

		})
	}
}
