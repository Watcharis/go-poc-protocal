package router

import (
	"context"
	"net/http"

	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/middleware"
	"watcharis/go-poc-protocal/restful_api/pooling/client/handlers"

	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/propagation"
)

func InitRouter(ctx context.Context, handler handlers.PoolingClientHandler) http.Handler {
	e := echo.New()
	e.GET("/health", echo.WrapHandler(http.HandlerFunc(pkg.HealthCheck)))

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // รองรับ W3C
		propagation.Baggage{},      // รองรับการส่งข้อมูล Metadata
	)

	e.Use(otelecho.Middleware(dto.APP_NAME, otelecho.WithPropagators(propagator)))
	e.Use(middleware.AddProjectNameFromContext(ctx))

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.POST("/pooling-server", handler.PoolingClient)

	return e
}
