package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/postgres"
)

// PaymentService exposes payment operations as Connect handlers.
type PaymentService interface {
	MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.MarkInvoicePaidRequest]) (*connect.Response[paymentv1.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[paymentv1.PayInvoiceRequest]) (*connect.Response[paymentv1.Invoice], error)
}

type paymentService struct {
	store postgres.DataStore
}

// NewPaymentService returns a PaymentService backed by the given DataStore.
func NewPaymentService(data postgres.DataStore) PaymentService {
	return &paymentService{store: data}
}

func (s *paymentService) MarkInvoicePaid(_ context.Context, _ *connect.Request[paymentv1.MarkInvoicePaidRequest]) (*connect.Response[paymentv1.Invoice], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("MarkInvoicePaid not implemented"))
}

func (s *paymentService) PayInvoice(_ context.Context, _ *connect.Request[paymentv1.PayInvoiceRequest]) (*connect.Response[paymentv1.Invoice], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("PayInvoice not implemented"))
}

func (s *paymentService) MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	return s.store.MakePayment(ctx, c)
}
