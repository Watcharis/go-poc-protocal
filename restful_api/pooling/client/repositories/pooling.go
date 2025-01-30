package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"watcharis/go-poc-protocal/pkg/httpclient"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/restful_api/pooling/client/models"

	"go.uber.org/zap"
)

type poolingRepository struct{}

func NewPoolingRepository() PoolingRepository {
	return &poolingRepository{}
}

func (r *poolingRepository) TrickerPoolingServer(ctx context.Context, request models.TrickerPoolingServerRequest) (models.TrickerPoolingServerResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return models.TrickerPoolingServerResponse{}, err
	}
	logger.Info(ctx, "Request body", zap.String("request", string(data)))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, models.TRICKER_POOLING_SEVER_URL, bytes.NewReader(data))
	if err != nil {
		return models.TrickerPoolingServerResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	httpClient := httpclient.CreateHttpClient()

	res, err := httpClient.Do(req)
	if err != nil {
		return models.TrickerPoolingServerResponse{}, err
	}

	defer res.Body.Close()

	responseByte, err := io.ReadAll(res.Body)
	if err != nil {
		return models.TrickerPoolingServerResponse{}, err
	}
	logger.Info(ctx, "Response body", zap.String("response", string(responseByte)))

	if res.StatusCode != http.StatusOK {
		return models.TrickerPoolingServerResponse{}, fmt.Errorf("http invalid status_code : %d", res.StatusCode)
	}

	var response models.TrickerPoolingServerResponse
	if err := json.Unmarshal(responseByte, &response); err != nil {
		return models.TrickerPoolingServerResponse{}, err
	}

	return response, nil
}

func (r *poolingRepository) GetStatus(ctx context.Context, request models.TrickerPoolingServerDataResponse) (models.GetStatusResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return models.GetStatusResponse{}, err
	}
	logger.Info(ctx, "Request body", zap.String("request", string(data)))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, models.GET_STATUS_POOLING_SEVER_URL, bytes.NewReader(data))
	if err != nil {
		return models.GetStatusResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	httpClient := httpclient.CreateHttpClient()

	res, err := httpClient.Do(req)
	if err != nil {
		return models.GetStatusResponse{}, err
	}

	defer res.Body.Close()

	responseByte, err := io.ReadAll(res.Body)
	if err != nil {
		return models.GetStatusResponse{}, err
	}
	logger.Info(ctx, "Response body", zap.String("response", string(responseByte)))

	if res.StatusCode != http.StatusOK {
		return models.GetStatusResponse{}, fmt.Errorf("http invalid status_code : %d", res.StatusCode)
	}

	var response models.GetStatusResponse
	if err := json.Unmarshal(responseByte, &response); err != nil {
		return models.GetStatusResponse{}, err
	}

	return response, nil
}
