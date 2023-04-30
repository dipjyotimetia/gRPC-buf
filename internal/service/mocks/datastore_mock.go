package mocks

import (
	"context"

	"github.com/bufbuild/connect-go"
	payment "github.com/grpc-buf/internal/gen/payment"
	"github.com/grpc-buf/internal/gen/registration"
	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockDataStore) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockDataStore) MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*connect.Response[payment.PaymentResponse]), args.Error(1)
}
