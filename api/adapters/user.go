package api_adapters

import (
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
	userv1 "github.com/okpalaChidiebere/chirper-app-gen-protos/user/v1"
)

func UserToProto (t *model.User) *userv1.User{
	return &userv1.User{
		Id: t.Id,
		Name: t.Name,
		AvatarUrl: t.AvatarURL,
		Tweets: t.Tweets,
	}
}

func UsersToProto (ts []*model.User) []*userv1.User{
	var users []*userv1.User
	for _, t := range ts {
		users = append(users, UserToProto(t))
	}
	return users
}
