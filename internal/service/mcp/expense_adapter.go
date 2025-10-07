package mcp

import (
	"context"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExpenseServiceAdapter adapts Connect-based ExpenseService to MCP interface
type ExpenseServiceAdapter struct {
	svc service.ExpenseService
}

// NewExpenseServiceAdapter creates a new adapter
func NewExpenseServiceAdapter(svc service.ExpenseService) *ExpenseServiceAdapter {
	return &ExpenseServiceAdapter{svc: svc}
}

// CreateExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) CreateExpense(ctx context.Context, req *expensev1.CreateExpenseRequest) (*expensev1.Expense, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.CreateExpense(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// GetExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) GetExpense(ctx context.Context, req *expensev1.GetExpenseRequest) (*expensev1.Expense, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.GetExpense(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// ListExpenses adapts from MCP to Connect
func (a *ExpenseServiceAdapter) ListExpenses(ctx context.Context, req *expensev1.ListExpensesRequest) (*expensev1.ListExpensesResponse, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.ListExpenses(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// UpdateExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) UpdateExpense(ctx context.Context, req *expensev1.UpdateExpenseRequest) (*expensev1.Expense, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.UpdateExpense(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// DeleteExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) DeleteExpense(ctx context.Context, req *expensev1.DeleteExpenseRequest) (*timestamppb.Timestamp, error) {
	connectReq := connect.NewRequest(req)
	resp, err := a.svc.DeleteExpense(ctx, connectReq)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}
