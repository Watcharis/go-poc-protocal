package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	mockServices "watcharis/go-poc-protocal/restful_api/ratelimit/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_restFulAPIHandlers_VerifyOtpRatelimit(t *testing.T) {

	type fields struct {
		services *mockServices.MockServices
	}

	type args struct {
		ctx context.Context
		req models.VerifyOtpRatelimitRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		fields  fields
		args    args
		senario func(args args, fields fields)
		want    models.VerifyOtpRatelimitResponse
	}{
		// TODO: Add test cases.
		{
			name: "verify otp ratelimit success",
			fields: fields{
				services: mockServices.NewMockServices(ctrl),
			},
			args: args{
				ctx: context.Background(),
				req: models.VerifyOtpRatelimitRequest{
					Uuid: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args, fields fields) {

				fields.services.EXPECT().VerifyOtpRatelimit(args.ctx, args.req).Return(models.VerifyOtpRatelimitResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
					Data: &models.VerifyOtpRatelimitDataResponse{
						Otp: "200139",
					},
					Error: nil,
				}, nil)

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

			tt.senario(tt.args, tt.fields)

			handler := NewRestFulAPIHandlers(tt.fields.services)

			reqBytes, err := json.Marshal(tt.args.req)
			if err != nil {
				t.Errorf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequestWithContext(tt.args.ctx, http.MethodPost, "/api/v1/verify-otp-ratelimit", bytes.NewReader(reqBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.VerifyOtpRatelimit(tt.args.ctx)(rec, req)

			expectedResp := tt.want

			var actualResp models.VerifyOtpRatelimitResponse
			if err := json.NewDecoder(rec.Body).Decode(&actualResp); err != nil {
				t.Errorf("failed to marshal request body: %v", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NoError(t, err)
			assert.Equal(t, expectedResp, actualResp)
		})
	}
}
