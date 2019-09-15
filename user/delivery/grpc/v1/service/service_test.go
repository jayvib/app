// +build unit

package service

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/jayvib/clean-architecture/model"
	"github.com/jayvib/clean-architecture/user/delivery/grpc/v1/proto"
	"github.com/jayvib/clean-architecture/user/mocks"
	"github.com/jayvib/clean-architecture/utils/generateutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUserSearchServiceServer_Store(t *testing.T) {
	user := &model.User{
		ID:        generateutil.GenerateID(model.GetUserTableName()),
		Firstname: "Luffy",
		Lastname:  "Monkey",
		Email:     "luffy.monkey@gmail.com",
		Username:  "luffy.monkey",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdAt, _ := ptypes.TimestampProto(user.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(user.UpdatedAt)

	esmock := new(mocks.SearchEngine)
	svc := NewUserSearchService(esmock)
	ctx := context.Background()
	type args struct {
		ctx context.Context
		req *proto.StoreRequest
	}

	tests := []struct {
		name    string
		svc     proto.UserSearchServiceServer
		args    args
		mock    func()
		want    *proto.StoreResponse
		wantErr bool
	}{
		{
			name: "Successful",
			svc:  svc,
			args: args{
				ctx: ctx,
				req: &proto.StoreRequest{
					User: &proto.User{
						Id:        user.ID,
						Firstname: user.Firstname,
						Lastname:  user.Lastname,
						Email:     user.Email,
						Username:  user.Username,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
					Api: "v1",
				},
			},
			mock: func() {
				esmock.On("Store", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil).Once()
			},
			want: &proto.StoreResponse{
				Api: "v1",
				User: &proto.User{
					Id:        user.ID,
					Firstname: user.Firstname,
					Lastname:  user.Lastname,
					Email:     user.Email,
					Username:  user.Username,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			wantErr: false,
		},
		{
			name: "API Request is not v1",
			svc:  svc,
			args: args{
				ctx: ctx,
				req: &proto.StoreRequest{
					User: &proto.User{},
					Api:  "v2",
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty User",
			svc:  svc,
			args: args{
				ctx: ctx,
				req: &proto.StoreRequest{
					User: nil,
					Api:  "v2",
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.svc.Store(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, got, tt.want, "unmatched")
		})
	}
}
