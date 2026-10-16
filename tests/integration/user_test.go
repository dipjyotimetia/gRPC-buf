//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
	"github.com/stretchr/testify/require"
)

func TestUsersRegister(t *testing.T) {
	client := userv1connect.NewUserServiceClient(connect.NewClient(connecthttp.NewTransport(
		http.DefaultClient,
		"http://localhost:8080",
		connecthttp.WithGRPCWeb(),
	)))
	req := &userv1.RegisterRequest{
		Email:     "test@example.com",
		Password:  "Password1",
		FirstName: "Test",
		LastName:  "Last",
	}
	res, err := client.RegisterUser(context.Background(), req)
	require.NoError(t, err)
	t.Log(res)
}

func TestUsersLogin(t *testing.T) {
	client := userv1connect.NewUserServiceClient(connect.NewClient(connecthttp.NewTransport(
		http.DefaultClient,
		"http://localhost:8080",
		connecthttp.WithGRPC(),
	)))
	req := &userv1.LoginRequest{
		Email:    "test@example.com",
		Password: "Password1",
	}
	res, err := client.LoginUser(context.Background(), req)
	require.NoError(t, err)
	t.Log(res)
}
