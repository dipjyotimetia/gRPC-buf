package httptransport

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/grpc-buf/internal/gen/proto/expense/expensev1connect"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
	"github.com/grpc-buf/internal/service"
)

// NewMux wires RPC handlers and returns an http.ServeMux.
func NewMux(payment service.PaymentService, user service.UserService, expense service.ExpenseService) *http.ServeMux {
	return NewMuxWithInterceptors(payment, user, expense)
}

// NewMuxWithInterceptors wires RPC handlers with optional unary interceptors
// (e.g. rate limiting, auth).
func NewMuxWithInterceptors(
	payment service.PaymentService,
	user service.UserService,
	expense service.ExpenseService,
	interceptors ...connect.Interceptor,
) *http.ServeMux {
	mux := http.NewServeMux()
	compress1KB := connect.WithCompressMinBytes(1024)
	opts := []connect.HandlerOption{compress1KB}
	if len(interceptors) > 0 {
		opts = append(opts, connect.WithInterceptors(interceptors...))
	}
	mux.Handle(paymentv1connect.NewPaymentServiceHandler(payment, opts...))
	mux.Handle(expensev1connect.NewExpenseServiceHandler(expense, opts...))
	mux.Handle(userv1connect.NewUserServiceHandler(user, opts...))

	checker := grpchealth.NewStaticChecker(
		paymentv1connect.PaymentServiceName,
		expensev1connect.ExpenseServiceName,
		userv1connect.UserServiceName,
	)
	mux.Handle(grpchealth.NewHandler(checker, compress1KB))

	reflector := grpcreflect.NewStaticReflector(
		paymentv1connect.PaymentServiceName,
		expensev1connect.ExpenseServiceName,
		userv1connect.UserServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector, compress1KB))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, compress1KB))

	mux.HandleFunc("/livez", livezHandler)
	return mux
}

func livezHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Debug("livez write failed", "error", err)
	}
}
