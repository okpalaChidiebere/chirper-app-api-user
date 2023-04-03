package usersservice

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	usersrepo "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/data_access"
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
)

func Test_CreateUserToDynamoDb(t *testing.T) {
	testCases := []struct {
		name          string
		user *model.User

		repoError error
		expectedError error
	}{
		{
			name: "should return error with repo error",
			user: &model.User{Name: "Chidi", Id: "someID"},
			repoError: errors.New("repo error"),
			expectedError: errors.New("repo error"),
		},
		{
			name: "should return no error with repo doesn't error",
			user: &model.User{Name: "Chidi", Id: "someID"},
			repoError: nil,
			expectedError: nil,
		},
		{
			name: "should return error with userID not provided",
			user: &model.User{Name: "Chidi"},
			repoError: nil,
			expectedError: errors.New("id is required"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)

			repoMock := usersrepo.NewMockRepository(ctrl)

			// tc.user.Id = uuid.NewString()
			expectedRepoCall := 1

			if(tc.user.Id == ""){
				expectedRepoCall = 0
			}

			repoMock.EXPECT().CreateUserToDynamoDb(ctx, tc.user).Times(expectedRepoCall).Return(tc.repoError)

			service := New(repoMock)
			err := service.CreateUserToDynamoDb(ctx, tc.user)

			assert.Equal(t, tc.expectedError, err)

		})
	}
}

func Test_GetUsersFromDynamoDb(t *testing.T) {
	testCases := []struct {
		name          string
		limit int32

		repoError error
		expectedError error
	}{
		{
			name: "should error if repo error",
			limit: 25,
			repoError: errors.New("repo error"),
			expectedError: errors.New("repo error"),
		},
		{
			name: "should default to 10 limit if none provided",
			limit: 0,
			repoError: nil,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)

			repoMock := usersrepo.NewMockRepository(ctrl)

			var expectLimit interface{}
			expectLimit = tc.limit

			//we expect 10 of int32 type to be passed as params to the GetUsers repo method if limit passed to the service layer is lessThan/equalTo 0 in the mock call
			if (tc.limit <= 0){
				expectLimit = gomock.Eq(int32(10))
			}

			repoMock.EXPECT().GetUsersFromDynamoDb(ctx, expectLimit, "").Times(1).Return(nil, "", tc.repoError)

			service := New(repoMock)
			_, _, err := service.GetUsers(ctx, tc.limit, "")

			assert.Equal(t, tc.expectedError, err)
		})
	}
}