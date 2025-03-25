package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	mockServices "watcharis/go-poc-protocal/restful_api/ratelimit/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {

	ctx := context.WithValue(context.Background(), "", "")

	tp, err := trace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	// tracer := otel.Tracer(logger.APP_NAME)
	// ctx, span := tracer.Start(ctx, logger.PROJECT_RATELIMIT)
	// defer span.End()

	// Create logger with TraceID and SpanID automatically included
	logger.InitOtelZapLogger("develop")
	defer logger.Sync()
}

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
			name: "case 1",
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

			reqBytes, _ := json.Marshal(tt.args.req)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/verify-otp-ratelimit", bytes.NewReader(reqBytes))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.VerifyOtpRatelimit(tt.args.ctx)(rec, req)

			expectedResp := tt.want

			var actualResp models.VerifyOtpRatelimitResponse
			err := json.NewDecoder(rec.Body).Decode(&actualResp)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.NoError(t, err)
			assert.Equal(t, expectedResp, actualResp)
		})
	}
}
