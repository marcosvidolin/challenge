package main

import (
	"api/cmd/gin/middleware"
	"api/cmd/gin/router"
	"api/internal/adapter"
	"api/internal/service"
	"api/internal/usecase"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx = context.Background()

		// API configs
		apiPort = os.Getenv("API_PORT")

		// 32 bytes hex key
		cryptorKey = os.Getenv("CRYPTOR_KEY")

		// Database config
		dbHost   = os.Getenv("DB_HOST")
		dbPort   = os.Getenv("DB_PORT")
		dbUser   = os.Getenv("DB_USER")
		dbPass   = os.Getenv("DB_PASS")
		dbName   = os.Getenv("DB_NAME")
		dbDriver = os.Getenv("DB_DRIVER")

		// Cache config
		cacheURL  = os.Getenv("CACHE_URL")
		cachePass = os.Getenv("CACHE_PASS")

		// AMQP config
		amqpURL   = os.Getenv("AMQP_URL")
		amqpQueue = os.Getenv("AMQP_QUEUE")
	)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open(dbDriver, datasource)
	if err != nil {
		log.Fatalf("error when try to open a database conection: %v", err.Error())
	}
	defer db.Close()

	redisAdapter := adapter.NewRedisCache(ctx, cacheURL, cachePass, 0)
	defer redisAdapter.Close()

	postgresAdapter := adapter.NewPostgreAdapter(db)
	rabbitmqAdapter, err := adapter.NewRabbitMQAdapter(amqpURL, amqpQueue)
	if err != nil {
		log.Fatalf("error when try to open a mqp conection: %v", err.Error())
	}

	cryptor, err := service.NewUserCryptor(cryptorKey)
	if err != nil {
		log.Fatalf("error when try to create cryptor service: %v", err.Error())
	}

	upsertUsecase := usecase.NewUpsertUsecase(rabbitmqAdapter, postgresAdapter, redisAdapter)
	getByIDUsecase := usecase.NewGetByIDUsecase(postgresAdapter, redisAdapter, cryptor)
	searchUsecase := usecase.NewSearchUsecase(postgresAdapter, redisAdapter, cryptor)

	go func() {
		for {
			if err := upsertUsecase.Execute(ctx); err != nil {
				log.Fatalf("error when try to open the message broker conection: %v", err.Error())
			}
		}
	}()

	server := gin.New()
	server.Use(middleware.Error())
	api := server.Group("/api")

	healthRouter := router.NewHealthRouter(db, redisAdapter)
	healthRouter.HealthRouter(api)

	getByIdRouter := router.NewGetByIDRouter(getByIDUsecase)
	getByIdRouter.GetByIDRouter(api)

	searchRouter := router.NewSearchRouter(searchUsecase)
	searchRouter.SearchRouter(api)

	if err := server.Run(apiPort); err != nil {
		log.Fatalf("error when try to run the api: %v", err.Error())
	}
}
