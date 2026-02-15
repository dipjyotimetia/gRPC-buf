//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"connectrpc.com/connect"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
)

// checkServerAvailable verifies if the HTTP server is available
func checkServerAvailable(t *testing.T, url string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Server not available at %s - start the server to run this test", url)
		}
		t.Skipf("Cannot reach server at %s: %v", url, err)
	}
	if resp.Body != nil {
		resp.Body.Close()
	}
}

func TestUsersRegister(t *testing.T) {
	checkServerAvailable(t, "http://localhost:8080")
	
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
		return
	}
	if res != nil && res.Msg != nil {
		t.Log(res.Msg)
	}
}

func TestUsersLogin(t *testing.T) {
	checkServerAvailable(t, "http://localhost:8080")
	
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
		return
	}
	if res != nil && res.Msg != nil {
		t.Log(res.Msg)
	}
}
