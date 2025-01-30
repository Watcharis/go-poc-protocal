package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/pooling/server/handlers"
	"watcharis/go-poc-protocal/restful_api/pooling/server/repositories/cache"
	"watcharis/go-poc-protocal/restful_api/pooling/server/routers"
	"watcharis/go-poc-protocal/restful_api/pooling/server/services"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	PORT = ":8781"

	REDIS_HOST     = "localhost"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
	REDIS_DATABASE = 0
)

func main() {
	ctx := context.WithValue(context.Background(), dto.APP_NAME, dto.PROJECT_POOLING_SEVER)

	tp, err := trace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	// Create logger with TraceID and SpanID automatically included
	logger.InitOtelZapLogger("develop")
	defer logger.Sync()

	redisClient := initRedis(ctx)
	redisRepository := cache.NewRedisRepository(redisClient)

	poolingService := services.NewPoolingServerService(redisRepository)
	poolingHandlers := handlers.NewPoolingServerHandlers(poolingService)

	handler := routers.InitRouter(ctx, poolingHandlers)

	s := &http.Server{
		Addr:    PORT,
		Handler: handler,
		//ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}
	logger.Info(ctx, "Server runnig on http://localhost"+s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func initRedis(ctx context.Context) *redis.Client {
	redisAddr := net.JoinHostPort(REDIS_HOST, REDIS_PORT)
	redisClient := redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		DB:              0,
		MinIdleConns:    30,
		PoolSize:        30,
		ConnMaxIdleTime: 30,
		PoolTimeout:     time.Duration(30 * time.Second),
		Password:        "",
	})

	redisPing, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Panic(ctx, "redis ping failed", zap.Error(err))
	} else {
		// log.Printf("Redis ping response=%s", redisPing)
		logger.Info(ctx, "Redis ping response="+redisPing)
	}

	return redisClient
}
