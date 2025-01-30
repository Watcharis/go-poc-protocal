package handlers

import (
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/pooling/server/models"
	"watcharis/go-poc-protocal/restful_api/pooling/server/services"

	"github.com/labstack/echo/v4"
)

type PoolingServerHandlers interface {
	ServerRecieveTricker(c echo.Context) error
	GetStatus(c echo.Context) error
}

type poolingServerHandlers struct {
	poolingServerService services.PoolingServerService
}

func NewPoolingServerHandlers(poolingServerService services.PoolingServerService) PoolingServerHandlers {
	return &poolingServerHandlers{
		poolingServerService: poolingServerService,
	}
}

func (h *poolingServerHandlers) ServerRecieveTricker(c echo.Context) error {
	ctx := c.Request().Context()

	var request models.TrickerPoolingServerRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.TrickerPoolingServerResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.TrickerPoolingServerResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	result, err := h.poolingServerService.ServerRecieveTricker(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.TrickerPoolingServerResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *poolingServerHandlers) GetStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var request models.GetStatusRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.GetStatusResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.GetStatusResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	result, err := h.poolingServerService.GetStatus(ctx, request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.GetStatusResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, result)
}
