package services

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/pooling/server/models"
	"watcharis/go-poc-protocal/restful_api/pooling/server/repositories"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

const (
	REDISKEY_FLAG_POOLING = "SF|%s|%s"
)

type PoolingServerService interface {
	ServerRecieveTricker(ctx context.Context, request models.TrickerPoolingServerRequest) (models.TrickerPoolingServerResponse, error)
	GetStatus(ctx context.Context, request models.GetStatusRequest) (models.GetStatusResponse, error)
}

type poolingServerService struct {
	redisRepository repositories.RedisRepository
}

func NewPoolingServerService(redisRepository repositories.RedisRepository) PoolingServerService {
	return &poolingServerService{
		redisRepository: redisRepository,
	}
}

func (s *poolingServerService) ServerRecieveTricker(ctx context.Context, request models.TrickerPoolingServerRequest) (models.TrickerPoolingServerResponse, error) {
	nctx := context.WithValue(context.Background(), dto.APP_NAME, ctx.Value(dto.APP_NAME))
	traceCtx := trace.SpanFromContext(ctx).SpanContext()
	nctx = trace.ContextWithSpanContext(nctx, traceCtx)

	refID := pkg.GenerateUUID()
	go func(ctx context.Context, request models.TrickerPoolingServerRequest, refID string) {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(1*time.Minute))
		defer func() {
			cancel()
			fmt.Println("stop backgroud process!!")
		}()
		for {
			select {
			case <-ctx.Done():
				logger.Error(ctx, "process time out, context done", zap.String("uuid", request.UserID), zap.String("ref_id", refID))
				return
			default:
				logger.Info(ctx, "start process pooling", zap.String("uuid", request.UserID), zap.String("ref_id", refID))

				// set default value false
				redisKey := fmt.Sprintf(REDISKEY_FLAG_POOLING, request.UserID, refID)
				_, err := s.redisRepository.Set(ctx, redisKey, "false", 0)
				if err != nil {
					logger.Error(ctx, "redis set failed", zap.String("uuid", request.UserID), zap.String("ref_id", refID),
						zap.String("redisKey", redisKey), zap.Error(err))
					return
				}

				// set sleep hold process for set redis
				time.Sleep(time.Duration(5 * time.Second))

				_, err = s.redisRepository.Set(ctx, redisKey, "true", 0)
				if err != nil {
					logger.Error(ctx, "redis set failed", zap.String("uuid", request.UserID), zap.String("ref_id", refID),
						zap.String("redisKey", redisKey), zap.Error(err))
					return
				}
				logger.Info(ctx, "set redis pooling success", zap.String("uuid", request.UserID), zap.String("ref_id", refID))
				return
			}
		}
	}(nctx, request, refID)

	return models.TrickerPoolingServerResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
		Data: &models.TrickerPoolingServerDataResponse{
			UserID: request.UserID,
			RefID:  refID,
		},
	}, nil
}

func (s *poolingServerService) GetStatus(ctx context.Context, request models.GetStatusRequest) (models.GetStatusResponse, error) {

	redisKey := fmt.Sprintf(REDISKEY_FLAG_POOLING, request.UserID, request.RefID)
	result, err := s.redisRepository.Get(ctx, redisKey)
	if err != nil {
		logger.Error(ctx, "cannot get redis pooling", zap.String("uuid", request.UserID), zap.String("ref_id", request.RefID),
			zap.String("redisKey", redisKey), zap.Error(err))
		return models.GetStatusResponse{}, err
	}

	flag, err := strconv.ParseBool(result)
	if err != nil {
		logger.Error(ctx, "cannot get redis pooling", zap.String("uuid", request.UserID), zap.String("ref_id", request.RefID),
			zap.String("result", result), zap.Error(err))
		return models.GetStatusResponse{}, err
	}

	return models.GetStatusResponse{
		CommonResponse: response.SetCommonResponse(response.STATUS_SUCCESS, http.StatusOK),
		Data: &models.GetStatusDataResponse{
			Flag:   flag,
			UserID: request.UserID,
			RefID:  request.RefID,
		},
	}, nil
}

// redis server
// - key  SF|<user_uuid>|<ref_id>
// - value flag (bool) ; true || false
// CLI : redis SET

// - key TX|REATELIMIT|<uuid>
// - value count (int64); 1
// CLI : redis INCR
