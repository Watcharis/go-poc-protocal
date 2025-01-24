package logger

import (
	"context"
	"fmt"

	"github.com/blendle/zapdriver"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

const (
	APP_NAME          string = "go-poc-protocal"
	PROJECT_RATELIMIT string = "ratelimit-service"
)

var logger *otelzap.Logger

func Sync() {
	logger.Sync()
}

// func addTrace(ctx context.Context, fields []zapcore.Field) []zapcore.Field {
// 	// ดึง Span จาก Context ที่มีอยู่
// 	span := trace.SpanFromContext(ctx)
// 	if span == nil {
// 		return fields // ถ้า Context ไม่มี Span, ไม่ต้องทำอะไร
// 	}

// 	spanContext := span.SpanContext()
// 	if !spanContext.IsValid() {
// 		fmt.Println("spanContext in Invalid")
// 		return fields // ถ้า SpanContext ไม่ Valid, ไม่ต้องทำอะไร
// 	}

// 	// ดึงค่า TraceID และ SpanID
// 	traceId := spanContext.TraceID().String()
// 	spanId := spanContext.SpanID().String()
// 	isSample := spanContext.TraceFlags().IsSampled()

// 	fmt.Println("traceId :", traceId)
// 	fmt.Println("spanId :", spanId)
// 	fmt.Println("isSample :", isSample)

// 	// เพิ่มข้อมูล traceID, spanID ลงใน Log Fields
// 	projectName, ok := ctx.Value("go-poc-protocal").(string)
// 	if ok {
// 		fields = append(fields, []zapcore.Field{
// 			zap.String("traceId", traceId),
// 			zap.String("spanId", spanId),
// 			zap.Bool("isSample", isSample),
// 			zap.String("projectName", projectName),
// 		}...)
// 	} else {
// 		fields = append(fields, []zapcore.Field{
// 			zap.String("traceId", traceId),
// 			zap.String("spanId", spanId),
// 			zap.Bool("isSample", isSample),
// 		}...)
// 	}

// 	return fields
// }

func addTrace(ctx context.Context, fields []zapcore.Field) []zapcore.Field {
	// Add traceID and spanID to logger fields using Context
	if span := trace.SpanFromContext(ctx); span != nil {

		spanContext := span.SpanContext()

		if spanContext.IsValid() {
			traceId := spanContext.TraceID().String()
			spanId := spanContext.SpanID().String()
			isSample := spanContext.TraceFlags().IsSampled()

			fmt.Println("traceId :", traceId)
			fmt.Println("spanId :", spanId)
			fmt.Println("isSample :", isSample)

			projectName, ok := ctx.Value(APP_NAME).(string)
			if ok {
				fields = append(fields, zapdriver.TraceContext(traceId, spanId, isSample, projectName)...)
			} else {
				fields = append(fields, zapdriver.TraceContext(traceId, spanId, isSample, "")...)
			}
		}
	}
	return fields
}

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Info(msg, addTrace(ctx, fields)...)
}

func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Warn(msg, addTrace(ctx, fields)...)
}

func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Error(msg, addTrace(ctx, fields)...)
}

func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Debug(msg, addTrace(ctx, fields)...)
}

func DPanic(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).DPanic(msg, addTrace(ctx, fields)...)
}

func Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Panic(msg, addTrace(ctx, fields)...)
}
