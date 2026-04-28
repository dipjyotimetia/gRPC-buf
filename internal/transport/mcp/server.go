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
	"github.com/grpc-buf/internal/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime/gosdk"
)

const serverName = "grpc-buf-mcp-server"

// Server exposes the Connect services as MCP tools over stdio. The underlying
// *mcp.Server is retained so Serve can drive its own transport loop.
type Server struct {
	raw *mcp.Server
}

// NewServer wires the Connect services into an MCP server using the official
// modelcontextprotocol/go-sdk via the protoc-gen-go-mcp gosdk adapter.
func NewServer(dataStore postgres.DataStore) (*Server, error) {
	raw, registrar := gosdk.NewServer(serverName, version.Get().Version)

	expenseSvc := service.NewExpenseService(dataStore)
	userSvc := service.NewUserService(dataStore)
	paymentSvc := service.NewPaymentService(dataStore)

	expenseAdapter := mcpadapter.NewExpenseServiceAdapter(expenseSvc)
	userAdapter := mcpadapter.NewUserServiceAdapter(userSvc)
	paymentAdapter := mcpadapter.NewPaymentServiceAdapter(paymentSvc)

	expensev1mcp.RegisterExpenseServiceHandler(registrar, expenseAdapter)
	userv1mcp.RegisterUserServiceHandler(registrar, userAdapter)
	paymentv1mcp.RegisterPaymentServiceHandler(registrar, paymentAdapter)

	slog.Info("MCP server initialized with all service handlers")

	return &Server{raw: raw}, nil
}

// Serve runs the MCP server over stdio until ctx is cancelled or the transport
// terminates.
func (s *Server) Serve(ctx context.Context) error {
	slog.Info("Starting MCP server with stdio transport")
	return s.raw.Run(ctx, &mcp.StdioTransport{})
}
