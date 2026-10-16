package httptransport

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/grpchealth/v2"
	"connectrpc.com/grpcreflect/v2"
	"github.com/grpc-buf/internal/gen/proto/expense/expensev1connect"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
	"github.com/grpc-buf/internal/service"
)

// NewMux wires RPC handlers and returns an http.ServeMux.
func NewMux(payment service.PaymentService, user service.UserService, expense service.ExpenseService) *http.ServeMux {
	return NewMuxWithInterceptors(payment, user, expense)
}

// NewMuxWithInterceptors wires RPC handlers with optional server interceptors
// (e.g. rate limiting, auth). Interceptors apply to the application services
// only; health checks and reflection are served without them.
func NewMuxWithInterceptors(
	payment service.PaymentService,
	user service.UserService,
	expense service.ExpenseService,
	interceptors ...connect.ServerInterceptor,
) *http.ServeMux {
	mux := http.NewServeMux()
	compress1KB := connecthttp.WithCompressMinBytes(1024)

	server := connect.NewServer(interceptors...)
	paymentv1connect.RegisterPaymentServiceHandler(server, payment)
	expensev1connect.RegisterExpenseServiceHandler(server, expense)
	userv1connect.RegisterUserServiceHandler(server, user)
	connecthttp.Mount(mux, server, compress1KB)

	serviceNames := []string{
		paymentv1connect.PaymentServiceName,
		expensev1connect.ExpenseServiceName,
		userv1connect.UserServiceName,
	}
	infra := connect.NewServer()
	grpchealth.Register(infra, grpchealth.NewStaticChecker(serviceNames...))
	grpcreflect.Register(infra, grpcreflect.WithNamer(grpcreflect.NamerFunc(func() []string {
		return serviceNames
	})))
	connecthttp.Mount(mux, infra, compress1KB)

	mux.HandleFunc("/livez", livezHandler)
	return mux
}

func livezHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Debug("livez write failed", "error", err)
	}
}
