package service

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres"
)

// UserService interface
type UserService interface {
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

type userService struct {
	userDatabase postgres.DataStore
}

func NewUserService(data postgres.DataStore) UserService {
	return &userService{userDatabase: data}
}

func (u userService) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	return u.userDatabase.LoginUser(ctx, req)
}

func (u userService) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	return u.userDatabase.RegisterUser(ctx, req)
}
