package mcp

import (
	"context"

	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
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
	return a.svc.CreateExpense(ctx, req)
}

// GetExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) GetExpense(ctx context.Context, req *expensev1.GetExpenseRequest) (*expensev1.Expense, error) {
	return a.svc.GetExpense(ctx, req)
}

// ListExpenses adapts from MCP to Connect
func (a *ExpenseServiceAdapter) ListExpenses(ctx context.Context, req *expensev1.ListExpensesRequest) (*expensev1.ListExpensesResponse, error) {
	return a.svc.ListExpenses(ctx, req)
}

// UpdateExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) UpdateExpense(ctx context.Context, req *expensev1.UpdateExpenseRequest) (*expensev1.Expense, error) {
	return a.svc.UpdateExpense(ctx, req)
}

// DeleteExpense adapts from MCP to Connect
func (a *ExpenseServiceAdapter) DeleteExpense(ctx context.Context, req *expensev1.DeleteExpenseRequest) (*emptypb.Empty, error) {
	return a.svc.DeleteExpense(ctx, req)
}
