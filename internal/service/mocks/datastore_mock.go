package mocks

import (
	"context"

	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/stretchr/testify/mock"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) LoginUser(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*userv1.LoginResponse), args.Error(1)
}

func (m *MockDataStore) RegisterUser(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*userv1.RegisterResponse), args.Error(1)
}

func (m *MockDataStore) MakePayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*paymentv1.PaymentResponse), args.Error(1)
}

func (m *MockDataStore) Close() {
	m.Called()
}
