package handlers

import "github.com/labstack/echo/v4"

type PoolingClientHandler interface {
	PoolingClient(c echo.Context) error
}
