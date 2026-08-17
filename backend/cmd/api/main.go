package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/saddiqs1/gymshark-saddiqs1/config"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/api"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/middleware"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/packsizes"
	"github.com/saddiqs1/gymshark-saddiqs1/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	appLogger := logger.New(cfg.Environment.LogLevel, cfg.Environment.AppEnv)
	dynamoClient := dynamodb.NewFromConfig(cfg.AwsConfig, func(options *dynamodb.Options) {
		if cfg.DynamoDB.Endpoint != "" {
			options.BaseEndpoint = &cfg.DynamoDB.Endpoint
		}
	})
	packSizeStore := packsizes.NewDynamoDBStore(dynamoClient, cfg.DynamoDB.TableName)
	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.WithRequestLogger(appLogger, api.NewRouter(appLogger, packSizeStore)),
	}

	appLogger.Info().Msgf("server running on %s", server.Addr)

	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		appLogger.Fatal().Err(err).Msg("server stopped")
	}
}
