package logger

import (
	"context"
	"fmt"
	"log"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		logger.Fatal(errMsg)
	}
	// Create trace provider without the exporter.
	sampler := sdktrace.ParentBased(
		sdktrace.AlwaysSample(),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
	)

	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		// sdktrace.WithResource(resource.Default()),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

// NewZapWithTracing creates a logger that supports adding traceID and spanID
func NewZapWithTracing() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(log.Writer()),
		zapcore.DebugLevel,
	)

	zapLogger := zap.New(core)

	logger = otelzap.New(zapLogger)
}
