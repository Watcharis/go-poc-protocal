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

func Test_restFulAPIHandlers_CreateOtp(t *testing.T) {
	type fields struct {
		services *mockServices.MockServices
	}

	type args struct {
		ctx context.Context
		req models.OtpRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				services: mockServices.NewMockServices(ctrl),
			},
			args: args{
				ctx: context.Background(),
				req: models.OtpRequest{
					UUID: "1234",
					Otp:  "200139",
				},
			},
			senario: func(args args, fields fields) {
				fields.services.EXPECT().CreateOtp(args.ctx, args.req).Return(models.OtpResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
					Error:          nil,
				}, nil)
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

			h := NewRestFulAPIHandlers(tt.fields.services)

			reqBytes, err := json.Marshal(tt.args.req)
			if err != nil {
				t.Errorf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequestWithContext(tt.args.ctx, http.MethodPost, "/api/v1/create-otp", bytes.NewReader(reqBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.CreateOtp(tt.args.ctx)(rec, req)

			expectedResp := tt.want

			var actualResp models.OtpResponse
			if err := json.NewDecoder(rec.Body).Decode(&actualResp); err != nil {
				t.Errorf("failed to marshal request body: %v", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NoError(t, err)
			assert.Equal(t, expectedResp, actualResp)
		})
	}
}
