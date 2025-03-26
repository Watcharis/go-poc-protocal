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

func Test_restFulAPIHandlers_CreateUserProfile(t *testing.T) {
	type fields struct {
		services *mockServices.MockServices
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
				services: mockServices.NewMockServices(ctrl),
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

				fields.services.EXPECT().CreateUserProfile(args.ctx, args.req).Return(models.ProifleResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				}, nil)

			},
			want: models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
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

			req := httptest.NewRequestWithContext(tt.args.ctx, http.MethodPost, "/api/v1/create-user-profile", bytes.NewReader(reqBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.CreateUserProfile(tt.args.ctx)(rec, req)

			expectedResp := tt.want

			var actualResp models.ProifleResponse
			if err := json.NewDecoder(rec.Body).Decode(&actualResp); err != nil {
				t.Errorf("failed to marshal request body: %v", err)
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NoError(t, err)
			assert.Equal(t, expectedResp, actualResp)
		})
	}
}

func Test_restFulAPIHandlers_GetUserProfile(t *testing.T) {
	type fields struct {
		services *mockServices.MockServices
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
		senario func(args args, fields fields)
		args    args
		want    models.ProifleResponse
	}{
		// TODO: Add test cases.
		{
			name: "get user profile success",
			fields: fields{
				services: mockServices.NewMockServices(ctrl),
			},
			args: args{
				ctx:  context.Background(),
				uuid: "1234",
			},
			senario: func(args args, fields fields) {

				fields.services.EXPECT().GetUserProfile(args.ctx, args.uuid).Return(models.ProifleResponse{
					CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
				}, nil)

			},
			want: models.ProifleResponse{
				CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.senario(tt.args, tt.fields)

			h := NewRestFulAPIHandlers(tt.fields.services)

			req := httptest.NewRequestWithContext(tt.args.ctx, http.MethodGet, "/api/v1/get-user-profile", nil)
			req.Header.Set("uuid", tt.args.uuid)

			rec := httptest.NewRecorder()

			h.GetUserProfile(tt.args.ctx)(rec, req)

			expectedResp := tt.want

			var actualResp models.ProifleResponse
			err := json.NewDecoder(rec.Body).Decode(&actualResp)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NoError(t, err)
			assert.Equal(t, expectedResp, actualResp)
		})
	}
}
