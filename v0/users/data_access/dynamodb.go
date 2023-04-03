package userdataaccess

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/okpalaChidiebere/chirper-app-api-user/v0/common"
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
)

//We can call this an Adapter! It connects to external service
type DynamoDbRepository struct {
	client common.DynamoDBAPI
	tableName string
}

type NextKey struct {
	Id string `json:"id"`
}

func NewDynamoDbRepo(client common.DynamoDBAPI, tableName string) *DynamoDbRepository{
	return &DynamoDbRepository{
		client: client,
		tableName: tableName,
	}
}

// Store creates a new user in the in users table.
func (r *DynamoDbRepository) CreateUserToDynamoDb(ctx context.Context, user *model.User) error {
	item, _ := attributevalue.MarshalMap(user)

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	}

	if len(user.Tweets) == 0{
		delete(item, "tweets")
	}

	// Write the new item to DynamoDB database
	if _, err := r.client.PutItem(ctx,input); err != nil {
		return  err
	}
	return nil
}

func (r *DynamoDbRepository) GetUsersFromDynamoDb(ctx context.Context, limit int32, nextKey string) ([]*model.User, string, error) {
	var input *dynamodb.ScanInput
	pItems := []*model.User{}

	if nextKey == "" {
		input = &dynamodb.ScanInput{
			TableName: aws.String(r.tableName),
			Limit:      aws.Int32(limit),
		}
	} else {
		nk := &NextKey{}

		//We decode the key
		k, _ := url.QueryUnescape(nextKey)

		//parse the key
		json.Unmarshal([]byte(k), nk)

		input = &dynamodb.ScanInput{
			TableName:  aws.String(r.tableName),
			Limit:      aws.Int32(limit),
			ExclusiveStartKey: map[string]types.AttributeValue{
				"id":        &types.AttributeValueMemberS{Value: nk.Id},
			},
		}
	}

	out, err := r.client.Scan(ctx, input)
	if err != nil {
		return pItems, "", err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &pItems)
	if err != nil {
		return pItems, "", err
	}

	var nxt NextKey
	if err := attributevalue.UnmarshalMap(out.LastEvaluatedKey, &nxt); err != nil {
		return pItems, "", err
	}

	var finalKeyValue string
	if nxt.Id == "" {
		//when the next key is null it means there is no more items ot return
		finalKeyValue = string("null")
	} else {
		out, _ := json.Marshal(nxt)
		finalKeyValue = url.QueryEscape(string(out))
	}

	return pItems, finalKeyValue, nil
}