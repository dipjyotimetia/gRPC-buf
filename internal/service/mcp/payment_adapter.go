package mcp

import (
	"context"

	"connectrpc.com/connect"
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
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.MakePayment(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// MarkInvoicePaid adapts from MCP to Connect
func (a *PaymentServiceAdapter) MarkInvoicePaid(ctx context.Context, req *paymentv1.Invoice) (*paymentv1.Invoice, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.MarkInvoicePaid(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// PayInvoice adapts from MCP to Connect
func (a *PaymentServiceAdapter) PayInvoice(ctx context.Context, req *paymentv1.Invoice) (*paymentv1.Invoice, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.PayInvoice(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}
