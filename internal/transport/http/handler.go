package httptransport

import (
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
    mux := http.NewServeMux()

	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(paymentv1connect.NewPaymentHandler(
		payment,
		compress1KB,
	))

	mux.Handle(expensev1connect.NewExpenseServiceHandler(
		expense,
		compress1KB,
	))

	mux.Handle(userv1connect.NewUserServiceHandler(
		user,
		compress1KB,
	))

	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			paymentv1connect.PaymentName,
			expensev1connect.ExpenseServiceName,
			userv1connect.UserServiceName),
		compress1KB,
	))

	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(
			paymentv1connect.PaymentName,
			expensev1connect.ExpenseServiceName,
			userv1connect.UserServiceName),
		compress1KB,
	))

	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(
			paymentv1connect.PaymentName,
			expensev1connect.ExpenseServiceName,
			userv1connect.UserServiceName,
		),
		compress1KB,
	))

	// Basic liveness endpoint for container runtime checks
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

    return mux
}

// NewMuxWithInterceptors allows custom unary interceptors (e.g., rate limiting).
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
    mux.Handle(paymentv1connect.NewPaymentHandler(payment, opts...))
    mux.Handle(expensev1connect.NewExpenseServiceHandler(expense, opts...))
    mux.Handle(userv1connect.NewUserServiceHandler(user, opts...))
    mux.Handle(grpchealth.NewHandler(
        grpchealth.NewStaticChecker(
            paymentv1connect.PaymentName,
            expensev1connect.ExpenseServiceName,
            userv1connect.UserServiceName),
        compress1KB,
    ))
    mux.Handle(grpcreflect.NewHandlerV1(
        grpcreflect.NewStaticReflector(
            paymentv1connect.PaymentName,
            expensev1connect.ExpenseServiceName,
            userv1connect.UserServiceName),
        compress1KB,
    ))
    mux.Handle(grpcreflect.NewHandlerV1Alpha(
        grpcreflect.NewStaticReflector(
            paymentv1connect.PaymentName,
            expensev1connect.ExpenseServiceName,
            userv1connect.UserServiceName,
        ),
        compress1KB,
    ))
    mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    return mux
}
