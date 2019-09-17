package service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/jayvib/app/model"
	"github.com/jayvib/app/user"
	"github.com/jayvib/app/user/delivery/grpc/v1/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotImplemented = errors.New("not implemented")

const apiVersion = "v1"

func NewUserSearchService(se user.SearchEngine) proto.UserSearchServiceServer {
	return &UserSearchServiceServer{
		searchEngine: se,
	}
}

// UserSearchServiceServer implements the proto.UserSearchService
type UserSearchServiceServer struct {
	searchEngine user.SearchEngine
}

func (u *UserSearchServiceServer) Update(ctx context.Context, in *proto.UpdateRequest) (response *proto.UpdateResponse, err error) {
	return nil, errNotImplemented
}

func (u *UserSearchServiceServer) Store(ctx context.Context, in *proto.StoreRequest) (response *proto.StoreResponse, err error) {
	if len(in.Api) > 0 {
		if in.Api != apiVersion {
			return nil, status.Error(codes.Unimplemented, "Unsupported API Version")
		}
	}

	if in.User == nil {
		return nil, status.Error(codes.InvalidArgument, "nil-value user")
	}
	userModel, err := userProtoToUserModel(in.User)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to convert userModel proto to userModel model")
	}
	err = u.searchEngine.Store(ctx, userModel)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to store")
	}

	return &proto.StoreResponse{Api: in.Api, User: in.User}, nil
}

func (u *UserSearchServiceServer) Delete(ctx context.Context, in *proto.DeleteRequest) (response *proto.DeleteResponse, err error) {
	return nil, errNotImplemented
}

func (u *UserSearchServiceServer) SearchByName(ctx context.Context, in *proto.SearchByNameRequest) (response *proto.SearchByNameResponse, err error) {
	return nil, errNotImplemented
}

func (u *UserSearchServiceServer) Search(ctx context.Context, in *proto.SearchRequest) (response *proto.SearchResponse, err error) {
	return nil, errNotImplemented
}

func userProtoToUserModel(u *proto.User) (*model.User, error) {
	user := &model.User{
		ID:        u.GetId(),
		Firstname: u.GetFirstname(),
		Lastname:  u.GetLastname(),
		Email:     u.GetEmail(),
		Username:  u.GetUsername(),
		Password:  u.GetPassword(),
	}
	createdAt, err := ptypes.Timestamp(u.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := ptypes.Timestamp(u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}

func userModelToUserProto(u *model.User) (*proto.User, error) {
	userProto := &proto.User{
		Id:        u.ID,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Lastname,
		Username:  u.Username,
		Password:  u.Password,
	}

	createdAt, err := ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	userProto.CreatedAt = createdAt
	userProto.UpdatedAt = updatedAt
	return userProto, nil
}
