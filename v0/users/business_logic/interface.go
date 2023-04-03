package usersservice

import (
	"context"

	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
)

//go:generate mockgen -destination mock.go -source=interface.go -package=usersservice
type Service interface {
	CreateUserToDynamoDb(ctx context.Context, user *model.User) error
	GetUsers(ctx context.Context, limit int32, nextKey string) ([]*model.User, string, error)
}