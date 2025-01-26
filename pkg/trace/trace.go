package trace

import (
	"context"
	"fmt"
	"net/http"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	// OpenTelemetry รองรับ Propagators หลายรูปแบบ เช่น W3C Trace Context หรือ B3 (ใช้กับ Jaeger และ Zipkin)
	//  การกำหนดค่า Propagator ให้เหมาะสมกับระบบ
	// 1.) ถ้าระบบของคุณใช้ Jaeger หรือ Zipkin:
	// 		- ใช้ propagation.B3{} หรือ jaeger.Jaeger{}
	// 		Jaeger Propagator:
	// 		Propagator ที่ใช้ใน Jaeger สำหรับ Context
	// 			Header:
	// 				uber-trace-id

	// 2.) ถ้าระบบใช้มาตรฐาน W3C:
	// 		- ใช้ propagation.TraceContext{} (Default)
	// 		W3C Trace Context (Default):
	// 		รองรับมาตรฐาน W3C Trace Context
	// 		ใช้ Trace ID และ Span ID ในการติดตาม Trace ระหว่าง Services
	// 		Header:
	// 			traceparent
	// 			tracestate

	// 3.) ถ้าต้องการส่ง Metadata เพิ่มเติม:
	// 		- ใช้ propagation.Baggage{} ควบคู่กับ Trace Context
	// 		Baggage Propagator:
	// 		ใช้สำหรับส่งข้อมูล Context เพิ่มเติม (Metadata) ระหว่าง Services
	//		Header:
	// 			baggage

	// 4.) ถ้าต้องการรองรับหลายมาตรฐานในเวลาเดียวกัน:
	// 		- ใช้ propagation.NewCompositeTextMapPropagator

	// example

	// propagator := propagation.NewCompositeTextMapPropagator(
	// 	// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
	// 	// takes precedence if both the traceparent and the XCTC headers exist.
	// 	propagator.CloudTraceOneWayPropagator{},
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// )

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // รองรับ W3C
		propagation.Baggage{},      // รองรับการส่งข้อมูล Metadata
	)

	otel.SetTextMapPropagator(propagator)

	return tp, nil
}

func MiddlewareAddTrace(ctx context.Context, next http.Handler) http.Handler {
	// ใช้ otelhttp.NewHandler Wrap http.Handle
	// เพื่อดึง Trace Context จาก Header
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectName := ctx.Value(dto.APP_NAME)
		ctx := context.WithValue(r.Context(), dto.APP_NAME, projectName)

		tracer := otel.Tracer(dto.APP_NAME)
		ctx, span := tracer.Start(ctx, dto.PROJECT_RATELIMIT,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}), dto.APP_NAME)
}
