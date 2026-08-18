//go:build integration

package packsizes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBStore(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	client := dynamodb.NewFromConfig(aws.Config{
		Region:      "local",
		Credentials: credentials.NewStaticCredentialsProvider("local", "local", ""),
	}, func(options *dynamodb.Options) {
		options.BaseEndpoint = aws.String(endpoint)
	})

	tableName := fmt.Sprintf("pack-sizes-integration-%d", time.Now().UnixNano())
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
	if err != nil {
		t.Fatalf("create integration test table at %s: %v", endpoint, err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if _, err := client.DeleteTable(cleanupCtx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}); err != nil {
			t.Errorf("delete integration test table: %v", err)
		}
	})

	if err := dynamodb.NewTableExistsWaiter(client).Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}, 10*time.Second); err != nil {
		t.Fatalf("wait for integration test table: %v", err)
	}

	store := NewDynamoDBStore(client, tableName)

	sizes, err := store.List(ctx)
	if err != nil {
		t.Fatalf("list empty store: %v", err)
	}
	if len(sizes) != 0 {
		t.Fatalf("empty store returned %v; expected no sizes", sizes)
	}

	for _, size := range []int{1000, 250, 500} {
		if err := store.Add(ctx, size); err != nil {
			t.Fatalf("add pack size %d: %v", size, err)
		}
	}

	sizes, err = store.List(ctx)
	if err != nil {
		t.Fatalf("list populated store: %v", err)
	}
	if expected := []int{250, 500, 1000}; !reflect.DeepEqual(sizes, expected) {
		t.Fatalf("List() returned %v; expected %v", sizes, expected)
	}

	if err := store.Add(ctx, 250); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("duplicate Add() returned %v; expected ErrAlreadyExists", err)
	}

	if err := store.Remove(ctx, 500); err != nil {
		t.Fatalf("remove existing pack size: %v", err)
	}
	sizes, err = store.List(ctx)
	if err != nil {
		t.Fatalf("list after removal: %v", err)
	}
	if expected := []int{250, 1000}; !reflect.DeepEqual(sizes, expected) {
		t.Fatalf("List() after removal returned %v; expected %v", sizes, expected)
	}

	if err := store.Remove(ctx, 500); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing Remove() returned %v; expected ErrNotFound", err)
	}
}
