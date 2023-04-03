package userdataaccess

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/okpalaChidiebere/chirper-app-api-user/v0/common"
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
	"github.com/stretchr/testify/assert"
)

type DynamodbMockClient struct {
	common.DynamoDBAPI
}

func (m DynamodbMockClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {

	result := dynamodb.PutItemOutput{}

	if params.Item == nil {
        return &result, errors.New("Missing required field PutItemInput.Item")
    }

    if params.TableName == nil || *params.TableName == "" {
        return &result, errors.New("Missing required field CreateTableInput.TableName")
    }

	return &result, nil
	
}

func (m DynamodbMockClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error){
	return &dynamodb.ScanOutput{
		Items:[]map[string]types.AttributeValue{
			{
				"id":        &types.AttributeValueMemberS{Value: "manu_ginobili"},
				"name":        &types.AttributeValueMemberS{Value: "Manu Ginobili"},
				"avatarUrl":        &types.AttributeValueMemberS{Value: ""},
				"tweets":        &types.AttributeValueMemberSS{Value: make([]string, 0)},
			},
		},
	}, nil
}

func initializeFakeDynamoDBRepository() (Repository, error) {
	return NewDynamoDbRepo(DynamodbMockClient{}, "fake-table-name"), nil
}

func Test_CreateUserToDynamoDb(t *testing.T) {

	testCases := []struct {
		name string

		//user to create
		user *model.User

		expectedError error
	}{
		{
			name: "Create should save user to storage",
			user: &model.User{},
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo, err := initializeFakeDynamoDBRepository()
			if err != nil {
				t.Fatalf("error initializing repository: %s", err.Error())
			}

			err = repo.CreateUserToDynamoDb(ctx, tc.user)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}

func Test_GetUsersFromDynamoDb(t *testing.T) {

	testCases := []struct {
		name string

		limit int32 
		nextKey string

		expectedResultLength int

		expectedError error
	}{
		{
			name: "GetUsers should return users pre created",
			expectedResultLength: 1,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo, err := initializeFakeDynamoDBRepository()
			if err != nil {
				t.Fatalf("error initializing repository: %s", err.Error())
			}

			users, _, err := repo.GetUsersFromDynamoDb(ctx, tc.limit, tc.nextKey)
			assert.Equal(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedResultLength, len(users))
		})
	}

}
//https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-profiles