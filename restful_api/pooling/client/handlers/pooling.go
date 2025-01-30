package handlers

import (
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/response"
	"watcharis/go-poc-protocal/restful_api/pooling/client/models"
	"watcharis/go-poc-protocal/restful_api/pooling/client/services"

	"github.com/labstack/echo/v4"
)

type poolingClientHandler struct {
	poolingClientService services.PoolingClientService
}

func NewPoolingClientHandler(poolingClientService services.PoolingClientService) PoolingClientHandler {
	return &poolingClientHandler{
		poolingClientService: poolingClientService,
	}
}

func (h *poolingClientHandler) PoolingClient(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.TrickerPoolingServerRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.PoolingClientResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	if err := pkg.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.PoolingClientResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusBadRequest),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	result, err := h.poolingClientService.ClientPoolingServer(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.PoolingClientResponse{
			CommonResponse: response.SetCommonResponse(response.STATUS_ERROR, http.StatusInternalServerError),
			Error: &response.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, result)
}
