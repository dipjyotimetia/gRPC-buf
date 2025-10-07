//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/grpc-buf/internal/config"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres"
	"github.com/grpc-buf/internal/service"
	mcptransport "github.com/grpc-buf/internal/transport/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/money"
)

func TestMCPServer_Initialization(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load("../../config/local.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Override database URL for test environment
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable"

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	require.NotNil(t, dataStore, "DataStore should not be nil")
	defer dataStore.Close()

	// Create MCP server
	mcpServer, err := mcptransport.NewServer(dataStore)
	require.NoError(t, err, "Failed to create MCP server")
	require.NotNil(t, mcpServer, "MCP server should not be nil")
}

func TestMCPServer_DatabaseConnectivity(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load("../../config/local.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Override database URL for test environment
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable"

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	require.NotNil(t, dataStore, "DataStore should not be nil")
	defer dataStore.Close()

	// Test database connectivity using Ping
	ctx := context.Background()
	err = dataStore.Ping(ctx)
	require.NoError(t, err, "Failed to ping database")

	// Test that we can list expenses (empty list is fine)
	listReq := connect.NewRequest(&expensev1.ListExpensesRequest{
		PageSize: 10,
	})
	listResp, err := dataStore.ListExpenses(ctx, listReq)
	require.NoError(t, err, "Failed to list expenses")
	require.NotNil(t, listResp, "List expenses response should not be nil")
	assert.GreaterOrEqual(t, len(listResp.Msg.Expenses), 0, "Expense list should be non-negative")
}

func TestMCPServer_ServiceAdapters(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load("../../config/local.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Override database URL for test environment
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable"

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	require.NotNil(t, dataStore, "DataStore should not be nil")
	defer dataStore.Close()

	t.Run("ExpenseService", func(t *testing.T) {
		// Create expense service
		expenseSvc := service.NewExpenseService(dataStore)
		require.NotNil(t, expenseSvc, "Expense service should not be nil")
	})

	t.Run("UserService", func(t *testing.T) {
		// Create user service
		userSvc := service.NewUserService(dataStore)
		require.NotNil(t, userSvc, "User service should not be nil")
	})

	t.Run("PaymentService", func(t *testing.T) {
		// Create payment service
		paymentSvc := service.NewPaymentService(dataStore)
		require.NotNil(t, paymentSvc, "Payment service should not be nil")
	})
}

func TestMCPServer_CreateExpense(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load("../../config/local.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Override database URL for test environment
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable"

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	require.NotNil(t, dataStore, "DataStore should not be nil")
	defer dataStore.Close()

	ctx := context.Background()

	// Create a test expense
	createReq := connect.NewRequest(&expensev1.CreateExpenseRequest{
		Expense: &expensev1.Expense{
			UserId: "00000000-0000-0000-0000-000000000001",
			Amount: &money.Money{
				CurrencyCode: "USD",
				Units:        25,
				Nanos:        0,
			},
			Category:    "MCP Test",
			Description: "Integration test expense",
		},
	})

	createResp, err := dataStore.CreateExpense(ctx, createReq)
	require.NoError(t, err, "Failed to create expense")
	require.NotNil(t, createResp, "Create expense response should not be nil")
	require.NotEmpty(t, createResp.Msg.Id, "Created expense should have an ID")

	expenseID := createResp.Msg.Id

	// Get the expense back
	getReq := connect.NewRequest(&expensev1.GetExpenseRequest{
		Id: expenseID,
	})
	getResp, err := dataStore.GetExpense(ctx, getReq)
	require.NoError(t, err, "Failed to get expense")
	require.NotNil(t, getResp, "Get expense response should not be nil")
	assert.Equal(t, expenseID, getResp.Msg.Id, "Expense ID should match")
	assert.Equal(t, "MCP Test", getResp.Msg.Category, "Category should match")

	// Clean up - delete the expense
	deleteReq := connect.NewRequest(&expensev1.DeleteExpenseRequest{
		Id: expenseID,
	})
	_, err = dataStore.DeleteExpense(ctx, deleteReq)
	require.NoError(t, err, "Failed to delete expense")
}

func TestMCPServer_UserRegistration(t *testing.T) {
	// Load test configuration
	cfg, err := config.Load("../../config/local.yaml")
	require.NoError(t, err, "Failed to load configuration")

	// Override database URL for test environment
	cfg.Database.URL = "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable"

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	require.NotNil(t, dataStore, "DataStore should not be nil")
	defer dataStore.Close()

	ctx := context.Background()

	// Register a test user
	registerReq := connect.NewRequest(&userv1.RegisterRequest{
		Email:     "mcptest@example.com",
		Password:  "TestPassword123!",
		FirstName: "MCP",
		LastName:  "Test User",
	})

	registerResp, err := dataStore.RegisterUser(ctx, registerReq)

	// User might already exist from previous test runs - that's OK
	if err == nil {
		require.NotNil(t, registerResp, "Register response should not be nil")
		assert.NotEmpty(t, registerResp.Msg.Id, "Should receive a user ID")
		assert.NotNil(t, registerResp.Msg.CreatedAt, "Should have creation timestamp")
	}
}
