package trace

import (
	"context"
	"fmt"
	"net/http"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// var tracer trace.Tracer

// setupTracer initializes the OpenTelemetry tracer
func SetupTracer(ctx context.Context, appName string) (*sdktrace.TracerProvider, error) {

	res, err := resource.New(ctx,
		resource.WithDetectors(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		errMsg := fmt.Sprintf("[error] init resource for tracer failed, err : %s", err.Error())
		logger.Fatal(ctx, errMsg)
	}
	// Create trace provider without the exporter.
	sampler := sdktrace.ParentBased(
		sdktrace.AlwaysSample(),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
	)

	// exporter, err := stdouttrace.New(
	// 	stdouttrace.WithPrettyPrint(),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	tp := sdktrace.NewTracerProvider(
		// sdktrace.WithBatcher(exporter),
		// sdktrace.WithResource(resource.Default()),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

func MiddlewareAddTrace(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tracer := otel.Tracer(dto.APP_NAME)
		ctx, span := tracer.Start(ctx, dto.PROJECT_RATELIMIT,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
