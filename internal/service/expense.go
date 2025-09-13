package service

import (
	"context"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ExpenseService interface
type ExpenseService interface {
	CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error)
	UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[timestamppb.Timestamp], error)
}

type expenseService struct {
	db postgres.DataStore
}

func NewExpenseService(data postgres.DataStore) ExpenseService {
	return &expenseService{db: data}
}

func (e expenseService) CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return e.db.CreateExpense(ctx, req)
}

func (e expenseService) GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return e.db.GetExpense(ctx, req)
}

func (e expenseService) ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error) {
	return e.db.ListExpenses(ctx, req)
}

func (e expenseService) UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return e.db.UpdateExpense(ctx, req)
}

func (e expenseService) DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[timestamppb.Timestamp], error) {
	return e.db.DeleteExpense(ctx, req)
}
