package integration

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
)

func TestUsersRegister(t *testing.T) {
	client := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPCWeb(),
	)
	req := connect.NewRequest(&userv1.RegisterRequest{
		Email:     "test@example.com",
		Password:  "Password1",
		FirstName: "Test",
		LastName:  "Last",
	})
	res, err := client.RegisterUser(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.Msg)
}

func TestUsersLogin(t *testing.T) {
	client := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)
	req := connect.NewRequest(&userv1.LoginRequest{
		Email:    "test@example.com",
		Password: "Password1",
	})
	res, err := client.LoginUser(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.Msg)
}
