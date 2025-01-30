package models

import "watcharis/go-poc-protocal/pkg/response"

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

type GetStatusRequest struct {
	UserID string `json:"user_id"`
	RefID  string `json:"ref_id"`
}

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
