package server

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/grpc-buf/internal/gen/proto/expense/expensev1connect"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
	"github.com/grpc-buf/internal/service"
)

func setupHandler(payment service.PaymentService, user service.UserService, expense service.ExpenseService) *http.ServeMux {
	mux := http.NewServeMux()

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		log.Fatal(err)
	}

	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(paymentv1connect.NewPaymentHandler(
		payment,
		compress1KB,
		connect.WithInterceptors(otelInterceptor),
	))

	mux.Handle(expensev1connect.NewExpenseServiceHandler(
		expense,
		compress1KB,
		connect.WithInterceptors(otelInterceptor),
	))

	mux.Handle(userv1connect.NewUserServiceHandler(
		user,
		compress1KB,
		connect.WithInterceptors(otelInterceptor),
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
