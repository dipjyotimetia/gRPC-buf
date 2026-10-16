package service

import (
	"context"

	"connectrpc.com/connect/v2"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/postgres"
)

// PaymentService exposes payment operations as Connect handlers.
type PaymentService interface {
	MakePayment(ctx context.Context, c *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error)
	MarkInvoicePaid(ctx context.Context, c *paymentv1.MarkInvoicePaidRequest) (*paymentv1.Invoice, error)
	PayInvoice(ctx context.Context, c *paymentv1.PayInvoiceRequest) (*paymentv1.Invoice, error)
}

type paymentService struct {
	store postgres.DataStore
}

// NewPaymentService returns a PaymentService backed by the given DataStore.
func NewPaymentService(data postgres.DataStore) PaymentService {
	return &paymentService{store: data}
}

func (s *paymentService) MarkInvoicePaid(_ context.Context, _ *paymentv1.MarkInvoicePaidRequest) (*paymentv1.Invoice, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, "MarkInvoicePaid not implemented")
}

func (s *paymentService) PayInvoice(_ context.Context, _ *paymentv1.PayInvoiceRequest) (*paymentv1.Invoice, error) {
	return nil, connect.NewError(connect.CodeUnimplemented, "PayInvoice not implemented")
}

func (s *paymentService) MakePayment(ctx context.Context, c *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	return s.store.MakePayment(ctx, c)
}
