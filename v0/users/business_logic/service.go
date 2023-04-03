package usersservice

import (
	"context"
	"errors"

	repo "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/data_access"
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
)

type ServiceImpl struct {
	repo repo.Repository
}


func New(repo repo.Repository) *ServiceImpl {
	return &ServiceImpl{repo}
}

func (s *ServiceImpl) CreateUserToDynamoDb(ctx context.Context, user *model.User) error{
	if user.Id == "" {
		return errors.New("id is required")
	}

	if err := s.repo.CreateUserToDynamoDb(ctx,user); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetUsers(ctx context.Context, limit int32, nextKey string) ([]*model.User, string, error){
	
	if (limit <= 0){
		limit = 10
	}
	return s.repo.GetUsersFromDynamoDb(ctx, limit, nextKey)
}