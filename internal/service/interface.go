package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentService interface
type PaymentService interface {
	MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
	PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error)
}

// UserService interface
type UserService interface {
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

// ExpenseService interface
type ExpenseService interface {
	CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error)
	UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[timestamppb.Timestamp], error)
}

type paymentService struct {
	paymentDatabase postgres.DataStore
}

type userService struct {
	userDatabase postgres.DataStore
}

type expenseService struct {
	db postgres.DataStore
}

func NewPaymentService(data postgres.DataStore) PaymentService {
	return &paymentService{
		paymentDatabase: data,
	}
}

func NewUserService(data postgres.DataStore) UserService {
	return &userService{
		userDatabase: data,
	}
}

func NewExpenseService(data postgres.DataStore) ExpenseService {
	return &expenseService{db: data}
}

func (u userService) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	return u.userDatabase.LoginUser(ctx, req)
}

func (u userService) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	return u.userDatabase.RegisterUser(ctx, req)
}

func (p paymentService) MarkInvoicePaid(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	// Return an explicit unimplemented error instead of panicking in production
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("MarkInvoicePaid not implemented"))
}

func (p paymentService) PayInvoice(ctx context.Context, c *connect.Request[paymentv1.Invoice]) (*connect.Response[paymentv1.Invoice], error) {
	// Return an explicit unimplemented error instead of panicking in production
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("PayInvoice not implemented"))
}

func (p paymentService) MakePayment(ctx context.Context, c *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	return p.paymentDatabase.MakePayment(ctx, c)
}

// Expense methods

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
	// Use timestamppb.Timestamp from protobuf in actual implementation; defer casting to datastore
	return e.db.DeleteExpense(ctx, req)
}
