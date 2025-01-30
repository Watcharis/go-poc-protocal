package services

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/pooling/client/models"
	"watcharis/go-poc-protocal/restful_api/pooling/client/repositories"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PoolingClientService interface {
	ClientPoolingServer(ctx context.Context, request models.TrickerPoolingServerRequest) (models.PoolingClientResponse, error)
}

type poolingClientService struct {
	poolingRepository repositories.PoolingRepository
}

func NewPoolinClientService(poolingRepository repositories.PoolingRepository) PoolingClientService {
	return &poolingClientService{
		poolingRepository: poolingRepository,
	}
}

func (s *poolingClientService) ClientPoolingServer(ctx context.Context, request models.TrickerPoolingServerRequest) (models.PoolingClientResponse, error) {

	responseTricker, err := s.poolingRepository.TrickerPoolingServer(ctx, request)
	if err != nil {
		logger.Error(ctx, "tricker pooling-server failed", zap.String("uuid", request.UserID), zap.Error(err))
		return models.PoolingClientResponse{}, err
	}

	if responseTricker.Code != http.StatusOK {
		err := fmt.Errorf("invalid pooling-server response_code : %d", responseTricker.Code)
		logger.Error(ctx, "invalid response_code", zap.Int("response_code", responseTricker.Code), zap.String("uuid", request.UserID), zap.Error(err))
		return models.PoolingClientResponse{}, err
	}

	if responseTricker.Data == nil {
		err := fmt.Errorf("invalid pooling-server empty response data : %v", responseTricker.Data)
		logger.Error(ctx, "Not found response data", zap.String("uuid", request.UserID), zap.Error(err))
		return models.PoolingClientResponse{}, err
	}

	nctx := context.WithValue(context.Background(), dto.APP_NAME, ctx.Value(dto.APP_NAME))
	traceCtx := trace.SpanFromContext(ctx).SpanContext()
	nctx = trace.ContextWithSpanContext(nctx, traceCtx)

	go func(ctx context.Context, res models.TrickerPoolingServerDataResponse) {
		// time.Sleep(time.Duration(5 * time.Second))
		logger.Info(ctx, "start process pooling to server", zap.String("uuid", res.UserID), zap.String("ref_id", res.RefID))

		ctx, cancel := context.WithTimeout(ctx, time.Duration(30*time.Second))
		defer func() {
			cancel()
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
				return
			default:
				// get status
				requestGetStatus := models.TrickerPoolingServerDataResponse{
					UserID: res.UserID,
					RefID:  res.RefID,
				}
				fmt.Printf("requestGetStatus : %+v\n", requestGetStatus)

				result, err := s.poolingRepository.GetStatus(ctx, requestGetStatus)
				if err != nil {
					logger.Error(ctx, "cannot get pooling status from server", zap.String("uuid", res.UserID),
						zap.String("ref_id", res.RefID), zap.Error(err))
					return
				}

				if result.Data.Flag {
					logger.Info(ctx, "get pooling status success", zap.Bool("flag", result.Data.Flag), zap.String("uuid", res.UserID),
						zap.String("ref_id", res.RefID), zap.Error(err))
					// do something
					return
				}

				time.Sleep(time.Duration(5 * time.Second))
			}
		}

	}(nctx, *responseTricker.Data)

	return models.PoolingClientResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
	}, nil
}
