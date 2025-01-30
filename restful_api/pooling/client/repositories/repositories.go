package repositories

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/pooling/client/models"
)

type PoolingRepository interface {
	TrickerPoolingServer(ctx context.Context, request models.TrickerPoolingServerRequest) (models.TrickerPoolingServerResponse, error)
	GetStatus(ctx context.Context, request models.TrickerPoolingServerDataResponse) (models.GetStatusResponse, error)
}
