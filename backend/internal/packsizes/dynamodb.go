package packsizes

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStore(client *dynamodb.Client, tableName string) *DynamoDBStore {
	return &DynamoDBStore{client: client, tableName: tableName}
}

func (s *DynamoDBStore) List(ctx context.Context) ([]int, error) {
	sizes := make([]int, 0)
	var startKey map[string]types.AttributeValue

	for {
		output, err := s.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                aws.String(s.tableName),
			ConsistentRead:           aws.Bool(true),
			ExclusiveStartKey:        startKey,
			ProjectionExpression:     aws.String("#size"),
			ExpressionAttributeNames: map[string]string{"#size": "size"},
		})
		if err != nil {
			return nil, fmt.Errorf("scan pack sizes: %w", err)
		}

		for _, item := range output.Items {
			attribute, ok := item["size"].(*types.AttributeValueMemberN)
			if !ok {
				return nil, fmt.Errorf("decode pack size: size is not a number")
			}
			size, err := strconv.Atoi(attribute.Value)
			if err != nil || size <= 0 {
				return nil, fmt.Errorf("decode pack size %q: expected a positive integer", attribute.Value)
			}
			sizes = append(sizes, size)
		}

		if len(output.LastEvaluatedKey) == 0 {
			break
		}
		startKey = output.LastEvaluatedKey
	}

	sort.Ints(sizes)
	return sizes, nil
}
