package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"watcharis/go-poc-protocal/pkg/consent"
	"watcharis/go-poc-protocal/pkg/dto"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/pkg/trace"
	"watcharis/go-poc-protocal/restful_api/ratelimit/handlers"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/cache"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories/db"
	"watcharis/go-poc-protocal/restful_api/ratelimit/router"
	"watcharis/go-poc-protocal/restful_api/ratelimit/services"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	SEVER_PORT = ":8777"

	REDIS_HOST     = "localhost"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
	REDIS_DATABASE = 0

	GORM_MYSQL_HOST     = "(localhost:3306)"
	GORM_MYSQL_USERNAME = "user"
	GORM_MYSQL_PASSWORD = "longpass"
	GORM_MYSQL_DB_NAME  = "lotto"
)

func main() {
	fmt.Println("start poc rest api")
	// ctx := context.Background()

	ctx := context.WithValue(context.Background(), dto.APP_NAME, dto.PROJECT_RATELIMIT)

	tp, err := trace.SetupTracer(ctx, dto.APP_NAME)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Panic(ctx, err.Error())
		}
	}()

	// tracer := otel.Tracer(logger.APP_NAME)
	// ctx, span := tracer.Start(ctx, logger.PROJECT_RATELIMIT)
	// defer span.End()

	// Create logger with TraceID and SpanID automatically included
	logger.InitOtelZapLogger("develop")
	defer logger.Sync()

	redisClient := initRedis(ctx)
	gormDB := initDatabase(ctx)

	consent.SetConsent(ctx, redisClient)

	redisRepository := cache.NewRedisRepository(redisClient)
	profileRepository := db.NewProfileRepository(gormDB)
	otpRepository := db.NewOtpRepository(gormDB)

	service := services.NewServices(redisRepository, profileRepository, otpRepository)

	handlers := handlers.NewRestFulAPIHandlers(service)

	routeHandlers := router.InitRouter(ctx, handlers)

	httpServer := http.Server{
		Addr:    SEVER_PORT,
		Handler: routeHandlers,
	}

	go func(port string) {
		defer httpServer.Close()
		// log.Printf("Server running on http://localhost%s\n", port)
		logger.Info(ctx, "Server runnig on http://localhost"+port)
		if err := httpServer.ListenAndServe(); err != nil {
			// log.Println("[error] cannot start server :", err)
			// logger.Error("cannot start server", zap.Error(err))
			logger.Panic(ctx, "cannot start server", zap.Error(err))
		}
	}(httpServer.Addr)

	wg := new(sync.WaitGroup)
	signal := make(chan os.Signal, 1)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		s := <-signal
		fmt.Println("signal :", s)
	}()
	wg.Wait()
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

func initDatabase(ctx context.Context) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
		GORM_MYSQL_USERNAME,
		GORM_MYSQL_PASSWORD,
		GORM_MYSQL_HOST,
		GORM_MYSQL_DB_NAME,
	)

	// log.Println("Initialing database with dsn")
	logger.Info(ctx, "Initialing database with dsn")

	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		logger.Panic(ctx, "gorm cannot connect msql", zap.Error(err))
	}

	// returns database statistics
	// log.Println("database is running")
	logger.Info(ctx, "database is running", zap.String("address", GORM_MYSQL_HOST))
	return db
}
