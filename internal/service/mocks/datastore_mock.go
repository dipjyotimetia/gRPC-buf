package mocks

import (
	"context"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*connect.Response[userv1.LoginResponse]), args.Error(1)
}

func (m *MockDataStore) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*connect.Response[userv1.RegisterResponse]), args.Error(1)
}

func (m *MockDataStore) MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*connect.Response[paymentv1.PaymentResponse]), args.Error(1)
}

func (m *MockDataStore) Close() {
	m.Called()
}
