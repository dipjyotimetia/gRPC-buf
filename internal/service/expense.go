package service

import (
	"context"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/postgres"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ExpenseService exposes expense CRUD as Connect handlers.
type ExpenseService interface {
	CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error)
	UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[emptypb.Empty], error)
}

type expenseService struct {
	store postgres.DataStore
}

// NewExpenseService returns an ExpenseService backed by the given DataStore.
func NewExpenseService(data postgres.DataStore) ExpenseService {
	return &expenseService{store: data}
}

func (s *expenseService) CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return s.store.CreateExpense(ctx, req)
}

func (s *expenseService) GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return s.store.GetExpense(ctx, req)
}

func (s *expenseService) ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error) {
	return s.store.ListExpenses(ctx, req)
}

func (s *expenseService) UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	return s.store.UpdateExpense(ctx, req)
}

func (s *expenseService) DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[emptypb.Empty], error) {
	return s.store.DeleteExpense(ctx, req)
}
