package models

import "watcharis/go-poc-protocal/pkg/response"

const (
	TRICKER_POOLING_SEVER_URL    = "http://localhost:8781/api/v1/tricker-pooling-server"
	GET_STATUS_POOLING_SEVER_URL = "http://localhost:8781/api/v1/get-status"
)

type PoolingClientResponse struct {
	response.CommonResponse
	Error *response.ErrorResponse `json:"error,omitempty"`
}

type TrickerPoolingServerRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type TrickerPoolingServerResponse struct {
	response.CommonResponse
	Data  *TrickerPoolingServerDataResponse `json:"data,omitempty"`
	Error *response.ErrorResponse           `json:"error,omitempty"`
}

type TrickerPoolingServerDataResponse struct {
	UserID string `json:"user_id"`
	RefID  string `json:"ref_id"`
}

// type GetStatusRequest struct {
// 	UserID string `json:"user_id"`
// 	RefID  string `json:"ref_id"`
// }

type GetStatusResponse struct {
	response.CommonResponse
	Data  *GetStatusDataResponse  `json:"data,omitempty"`
	Error *response.ErrorResponse `json:"error,omitempty"`
}

type GetStatusDataResponse struct {
	Flag   bool   `json:"flag"`
	UserID string `json:"user_id"`
	RefID  string `json:"ref_id"`
}

// redis server
// - key  F|<user_uuid>|<ref_id>
// - value flag (bool) ; true || false
// CLI : redis SET

// - key TX|REATELIMIT|<uuid>
// - value count (int64); 1
// CLI : redis INCR
