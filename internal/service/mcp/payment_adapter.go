package mcp

import (
	"context"

	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/service"
)

// PaymentServiceAdapter adapts Connect-based PaymentService to MCP interface
type PaymentServiceAdapter struct {
	svc service.PaymentService
}

// NewPaymentServiceAdapter creates a new adapter
func NewPaymentServiceAdapter(svc service.PaymentService) *PaymentServiceAdapter {
	return &PaymentServiceAdapter{svc: svc}
}

// MakePayment adapts from MCP to Connect
func (a *PaymentServiceAdapter) MakePayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	return a.svc.MakePayment(ctx, req)
}

// MarkInvoicePaid adapts from MCP to Connect
func (a *PaymentServiceAdapter) MarkInvoicePaid(ctx context.Context, req *paymentv1.MarkInvoicePaidRequest) (*paymentv1.Invoice, error) {
	return a.svc.MarkInvoicePaid(ctx, req)
}

// PayInvoice adapts from MCP to Connect
func (a *PaymentServiceAdapter) PayInvoice(ctx context.Context, req *paymentv1.PayInvoiceRequest) (*paymentv1.Invoice, error) {
	return a.svc.PayInvoice(ctx, req)
}
