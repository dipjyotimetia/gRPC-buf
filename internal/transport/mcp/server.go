package mcp

import (
	"context"
	"log/slog"

	expensev1mcp "github.com/grpc-buf/internal/gen/proto/expense/expensev1mcp"
	paymentv1mcp "github.com/grpc-buf/internal/gen/proto/payment/paymentv1mcp"
	userv1mcp "github.com/grpc-buf/internal/gen/proto/registration/userv1mcp"
	"github.com/grpc-buf/internal/postgres"
	"github.com/grpc-buf/internal/service"
	mcpadapter "github.com/grpc-buf/internal/service/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

// Server wraps the MCP server
type Server struct {
	mcp *mcpserver.MCPServer
}

// NewServer creates a new MCP server with all service handlers registered
func NewServer(dataStore postgres.DataStore) (*Server, error) {
	// Create MCP server with metadata
	mcpSrv := mcpserver.NewMCPServer(
		"grpc-buf-mcp-server",
		"1.0.0",
	)

	// Initialize services
	expenseSvc := service.NewExpenseService(dataStore)
	userSvc := service.NewUserService(dataStore)
	paymentSvc := service.NewPaymentService(dataStore)

	// Create adapters
	expenseAdapter := mcpadapter.NewExpenseServiceAdapter(expenseSvc)
	userAdapter := mcpadapter.NewUserServiceAdapter(userSvc)
	paymentAdapter := mcpadapter.NewPaymentServiceAdapter(paymentSvc)

	// Register handlers
	expensev1mcp.RegisterExpenseServiceHandler(mcpSrv, expenseAdapter)
	userv1mcp.RegisterUserServiceHandler(mcpSrv, userAdapter)
	paymentv1mcp.RegisterPaymentHandler(mcpSrv, paymentAdapter)

	slog.Info("MCP server initialized with all service handlers")

	return &Server{mcp: mcpSrv}, nil
}

// Serve starts the MCP server using stdio transport
func (s *Server) Serve(ctx context.Context) error {
	slog.Info("Starting MCP server with stdio transport")
	return mcpserver.ServeStdio(s.mcp)
}
