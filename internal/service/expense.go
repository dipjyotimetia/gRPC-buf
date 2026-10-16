package service

import (
	"context"

	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/postgres"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ExpenseService exposes expense CRUD as Connect handlers.
type ExpenseService interface {
	CreateExpense(ctx context.Context, req *expensev1.CreateExpenseRequest) (*expensev1.Expense, error)
	GetExpense(ctx context.Context, req *expensev1.GetExpenseRequest) (*expensev1.Expense, error)
	ListExpenses(ctx context.Context, req *expensev1.ListExpensesRequest) (*expensev1.ListExpensesResponse, error)
	UpdateExpense(ctx context.Context, req *expensev1.UpdateExpenseRequest) (*expensev1.Expense, error)
	DeleteExpense(ctx context.Context, req *expensev1.DeleteExpenseRequest) (*emptypb.Empty, error)
}

type expenseService struct {
	store postgres.DataStore
}

// NewExpenseService returns an ExpenseService backed by the given DataStore.
func NewExpenseService(data postgres.DataStore) ExpenseService {
	return &expenseService{store: data}
}

func (s *expenseService) CreateExpense(ctx context.Context, req *expensev1.CreateExpenseRequest) (*expensev1.Expense, error) {
	return s.store.CreateExpense(ctx, req)
}

func (s *expenseService) GetExpense(ctx context.Context, req *expensev1.GetExpenseRequest) (*expensev1.Expense, error) {
	return s.store.GetExpense(ctx, req)
}

func (s *expenseService) ListExpenses(ctx context.Context, req *expensev1.ListExpensesRequest) (*expensev1.ListExpensesResponse, error) {
	return s.store.ListExpenses(ctx, req)
}

func (s *expenseService) UpdateExpense(ctx context.Context, req *expensev1.UpdateExpenseRequest) (*expensev1.Expense, error) {
	return s.store.UpdateExpense(ctx, req)
}

func (s *expenseService) DeleteExpense(ctx context.Context, req *expensev1.DeleteExpenseRequest) (*emptypb.Empty, error) {
	return s.store.DeleteExpense(ctx, req)
}
