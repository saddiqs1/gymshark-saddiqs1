package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/saddiqs1/gymshark-saddiqs1/config"
	"github.com/saddiqs1/gymshark-saddiqs1/pkg/logger"
)

var initialPackSizes = []int{250, 500, 1000, 2000, 5000}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}
	appLogger := logger.New(cfg.Environment.LogLevel, cfg.Environment.AppEnv)
	if cfg.DynamoDB.TableName == "" {
		appLogger.Fatal().Msg("PACK_SIZES_TABLE_NAME must be configured")
	}

	client := dynamodb.NewFromConfig(aws.Config{
		Region:      "local",
		Credentials: credentials.NewStaticCredentialsProvider("local", "local", ""),
	}, func(options *dynamodb.Options) {
		if cfg.DynamoDB.Endpoint != "" {
			options.BaseEndpoint = aws.String(cfg.DynamoDB.Endpoint)
		}
	})

	if err := createTable(ctx, client, cfg.DynamoDB.TableName); err != nil {
		appLogger.Fatal().Err(err).Msg("error creating table")
	}

	waiter := dynamodb.NewTableExistsWaiter(client)
	if err := waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(cfg.DynamoDB.TableName),
	}, 20*time.Second); err != nil {
		appLogger.Fatal().Err(err).Msgf("wait for table %q: %v", cfg.DynamoDB.TableName, err)
	}

	for _, size := range initialPackSizes {
		if err := putPackSize(ctx, client, cfg.DynamoDB.TableName, size); err != nil {
			appLogger.Fatal().Err(err).Msgf("error inserting seed data: %v", size)
		}
	}

	appLogger.Info().Msgf("seeded table %q with pack sizes %v", cfg.DynamoDB.TableName, initialPackSizes)
}

func createTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("size"), AttributeType: types.ScalarAttributeTypeN},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("size"), KeyType: types.KeyTypeHash},
		},
	})
	var alreadyExists *types.ResourceInUseException
	if errors.As(err, &alreadyExists) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("create table %q: %w", tableName, err)
	}
	return nil
}

func putPackSize(ctx context.Context, client *dynamodb.Client, tableName string, size int) error {
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"size": &types.AttributeValueMemberN{Value: fmt.Sprint(size)},
		},
		ConditionExpression: aws.String("attribute_not_exists(#size)"),
		ExpressionAttributeNames: map[string]string{
			"#size": "size",
		},
	})
	var alreadyExists *types.ConditionalCheckFailedException
	if errors.As(err, &alreadyExists) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("seed pack size %d: %w", size, err)
	}
	return nil
}
