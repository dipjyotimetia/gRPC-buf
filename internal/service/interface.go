package service

import (
	"context"

	"connectrpc.com/connect"
	payment "github.com/grpc-buf/internal/gen/payment"
	userv1 "github.com/grpc-buf/internal/gen/registration"
	"github.com/grpc-buf/internal/mongo"
)

type paymentService struct {
	paymentDatabase mongo.DataStore
}

type userService struct {
	userDatabase mongo.DataStore
}

type PaymentService interface {
	MarkInvoicePaid(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error)
	MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error)
}

func NewPaymentService(data mongo.DataStore) PaymentService {
	return &paymentService{
		paymentDatabase: data,
	}
}

type UserService interface {
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

func NewUserService(data mongo.DataStore) UserService {
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

func (p paymentService) MarkInvoicePaid(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error) {
	// TODO implement me
	panic("implement me")
}

func (p paymentService) PayInvoice(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error) {
	// TODO implement me
	panic("implement me")
}

func (p paymentService) MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error) {
	return p.paymentDatabase.MakePayment(ctx, req)
}
