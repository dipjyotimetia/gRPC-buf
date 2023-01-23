package service

import (
	"context"

	"github.com/bufbuild/connect-go"
	payment "github.com/grpc-buf/internal/gen/payment"
	"github.com/grpc-buf/internal/mongo"
)

type paymentService struct {
	paymentDatabase mongo.DataStore
}

func NewPaymentService(data mongo.DataStore) PaymentService {
	return &paymentService{
		paymentDatabase: data,
	}
}

type PaymentService interface {
	MarkInvoicePaid(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[payment.Invoice]) (*connect.Response[payment.Invoice], error)
	MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error)
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
