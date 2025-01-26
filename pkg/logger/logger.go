package logger

import (
	"context"
	"log"
	"watcharis/go-poc-protocal/pkg/dto"

	"github.com/blendle/zapdriver"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *otelzap.Logger

func Sync() {
	logger.Sync()
}

// NewZapWithTracing creates a logger that supports adding traceID and spanID
func InitOtelZapLogger(env string) {

	var encoderCfg zapcore.EncoderConfig
	if env != "develop" {
		log.Println("init zap config in Env: " + env)
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		log.Println("init zap config in Env: " + env)
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(log.Writer()),
		zapcore.DebugLevel,
	)

	zapLogger := zap.New(core,
		zap.AddStacktrace(zap.InfoLevel),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddStacktrace(zap.WarnLevel),
		zap.AddStacktrace(zap.PanicLevel),
		zap.AddStacktrace(zap.DebugLevel),
		zap.AddStacktrace(zap.DPanicLevel),
		zap.AddStacktrace(zap.FatalLevel),
	)

	logger = otelzap.New(zapLogger)
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
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return fields
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return fields
	}

	traceId := spanContext.TraceID().String()
	spanId := spanContext.SpanID().String()
	isSample := spanContext.TraceFlags().IsSampled()

	projectName, ok := ctx.Value(dto.APP_NAME).(string)
	if ok {
		fields = append(fields, zapdriver.TraceContext(traceId, spanId, isSample, projectName)...)
	} else {
		fields = append(fields, zapdriver.TraceContext(traceId, spanId, isSample, "")...)
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

func Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	logger.Ctx(ctx).Fatal(msg, addTrace(ctx, fields)...)
}
