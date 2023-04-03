package userdataaccess

import (
	"context"

	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
)

/*
This interface is a Port. It provides an interface that is independent from specific technologies
A Port can have multiple Adapters. This means any adapter that belongs to this Port, will implement
these two methods. In our app, we just have one Adapter(UserDynamoDbRepository) that belongs to this Port (Repository)
*/
//go:generate mockgen -destination mock.go -source=interface.go -package=userdataaccess
type Repository interface {
	CreateUserToDynamoDb(ctx context.Context, user *model.User) error
	GetUsersFromDynamoDb(ctx context.Context, limit int32, nextKey string) ([]*model.User, string, error)
}
