package service

import (
	"context"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres"
)

// PaymentService interface
type PaymentService interface {
	MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
}

// UserService interface
type UserService interface {
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

type paymentService struct {
	paymentDatabase postgres.DataStore
}

type userService struct {
	userDatabase postgres.DataStore
}

func NewPaymentService(data postgres.DataStore) PaymentService {
	return &paymentService{
		paymentDatabase: data,
	}
}

func NewUserService(data postgres.DataStore) UserService {
	return &userService{
		userDatabase: data,
	}
}

func (u userService) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	return u.userDatabase.LoginUser(ctx, req)
}

func (u userService) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	return u.userDatabase.RegisterUser(ctx, req)
}

func (p paymentService) MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	// TODO implement me
	panic("implement me")
}

func (p paymentService) PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	// TODO implement me
	panic("implement me")
}

func (p paymentService) MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	return p.paymentDatabase.MakePayment(ctx, c)
}
