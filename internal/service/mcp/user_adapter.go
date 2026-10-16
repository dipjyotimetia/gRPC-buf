package mcp

import (
	"context"

	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/service"
)

// UserServiceAdapter adapts Connect-based UserService to MCP interface
type UserServiceAdapter struct {
	svc service.UserService
}

// NewUserServiceAdapter creates a new adapter
func NewUserServiceAdapter(svc service.UserService) *UserServiceAdapter {
	return &UserServiceAdapter{svc: svc}
}

// LoginUser adapts from MCP to Connect
func (a *UserServiceAdapter) LoginUser(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	return a.svc.LoginUser(ctx, req)
}

// RegisterUser adapts from MCP to Connect
func (a *UserServiceAdapter) RegisterUser(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	return a.svc.RegisterUser(ctx, req)
}
