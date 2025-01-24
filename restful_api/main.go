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
	"watcharis/go-poc-protocal/restful_api/handlers"
	"watcharis/go-poc-protocal/restful_api/repositories/cache"
	"watcharis/go-poc-protocal/restful_api/repositories/db"
	"watcharis/go-poc-protocal/restful_api/router"
	"watcharis/go-poc-protocal/restful_api/services"

	"github.com/redis/go-redis/v9"
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
	ctx := context.Background()

	redisClient := initRedis()
	gormDB := initDatabase()

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
		log.Printf("Server running on http://localhost%s\n", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("[error] cannot start server :", err)
			log.Panic(err)
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

func initRedis() *redis.Client {
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
		log.Panic(err)
	} else {
		log.Printf("Redis ping response=%s", redisPing)
	}

	return redisClient
}

func initDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
		GORM_MYSQL_USERNAME,
		GORM_MYSQL_PASSWORD,
		GORM_MYSQL_HOST,
		GORM_MYSQL_DB_NAME,
	)

	log.Println("Initialing database with dsn")

	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	// returns database statistics
	log.Println("database is running")
	return db
}
