package service

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres"
)

// UserService exposes registration and login as Connect handlers.
type UserService interface {
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

type userService struct {
	store postgres.DataStore
}

// NewUserService returns a UserService backed by the given DataStore.
func NewUserService(data postgres.DataStore) UserService {
	return &userService{store: data}
}

func (s *userService) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	return s.store.LoginUser(ctx, req)
}

func (s *userService) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	return s.store.RegisterUser(ctx, req)
}
