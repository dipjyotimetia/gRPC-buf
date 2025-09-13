package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/postgres"
)

// PaymentService interface
type PaymentService interface {
	MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
}

type paymentService struct {
	paymentDatabase postgres.DataStore
}

func NewPaymentService(data postgres.DataStore) PaymentService {
	return &paymentService{paymentDatabase: data}
}

func (p paymentService) MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("MarkInvoicePaid not implemented"))
}

func (p paymentService) PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("PayInvoice not implemented"))
}

func (p paymentService) MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	return p.paymentDatabase.MakePayment(ctx, c)
}
