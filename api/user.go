package api

import (
	"context"
	"log"

	connect "github.com/bufbuild/connect-go"
	apiadapters "github.com/okpalaChidiebere/chirper-app-api-user/api/adapters"
	usersservice "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/business_logic"
	model "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/model"
	imagefilterv1 "github.com/okpalaChidiebere/chirper-app-gen-protos/image_filter/v1"
	imagefilterconnect "github.com/okpalaChidiebere/chirper-app-gen-protos/image_filter/v1/image_filterv1connect"
	userv1 "github.com/okpalaChidiebere/chirper-app-gen-protos/user/v1"
	userv1connect "github.com/okpalaChidiebere/chirper-app-gen-protos/user/v1/userv1connect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)



type UserServer struct {
	// userv1connect.UnimplementedUserServiceHandler
	usersService usersservice.Service
	imageFilterServiceClient imagefilterconnect.ImagefilterServiceClient
}

func NewUserServer(usersService usersservice.Service, imageFilterServiceClient imagefilterconnect.ImagefilterServiceClient) userv1connect.UserServiceHandler {
	return &UserServer{ 
		usersService: usersService,
		imageFilterServiceClient: imageFilterServiceClient,
	}
}

func (s *UserServer) ListUsers(ctx context.Context, req *connect.Request[userv1.ListUsersRequest]) (*connect.Response[userv1.ListUsersResponse], error){
	users, nk, err := s.usersService.GetUsers(ctx, req.Msg.GetLimit(), req.Msg.GetNextKey())
	if err != nil {
		log.Printf("ListUsers Err: %v", err.Error())
		return nil, err
	}
	
	for i, p := range users {
		cr := connect.NewRequest(&imagefilterv1.GetGetSignedUrlRequest{
			ObjectKey: p.AvatarURL,
		})
		res, err := s.imageFilterServiceClient.GetGetSignedUrl(ctx, cr)
		if err != nil {
			log.Printf("imageFilterServiceClient.GetGetSignedUrl Err: %v", err.Error())
			continue
		}
		users[i].AvatarURL = res.Msg.GetUrl()
	}

	todos := apiadapters.UsersToProto(users)

	mRes := connect.NewResponse(&userv1.ListUsersResponse{
		Items: todos,
		NextKey: nk,
	})
	return mRes, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *connect.Request[userv1.CreateUserRequest]) (*connect.Response[emptypb.Empty], error){
	r := req.Msg.GetUser()

	//TODO: Connect to the Image filter grpc server to filter the image
	cr := connect.NewRequest(&imagefilterv1.FilterImageRequest{
		ImageUrl: r.AvatarUrl,
	})
	res, err := s.imageFilterServiceClient.FilterImage(ctx, cr)
	if err != nil {
		log.Printf("CreateUser Err: %s", err.Error())
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	item := model.User{
		Id: r.Id,
		Name: r.Name,
		AvatarURL: res.Msg.GetFilteredUrl(),
		Tweets: r.Tweets,
	}
	err = s.usersService.CreateUserToDynamoDb(ctx, &item)
	if err != nil {
		log.Printf("CreateUser Err: %s", err.Error())
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}


//https://github.com/jerryan999/book-service-rpc/blob/main/internal/server.go