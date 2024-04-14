package server

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/grpc-buf/internal/gen/payment/paymentv1connect"
	"github.com/grpc-buf/internal/gen/registration/userv1connect"
)

func setupHandler() *http.ServeMux {
	mux := http.NewServeMux()

	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(paymentv1connect.NewPaymentHandler(
		paymentService,
		compress1KB,
	))
	mux.Handle(userv1connect.NewUserServiceHandler(
		userService,
		compress1KB,
	))
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(
			paymentv1connect.PaymentName,
			userv1connect.UserServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(
			paymentv1connect.PaymentName,
			userv1connect.UserServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(
			paymentv1connect.PaymentName,
			userv1connect.UserServiceName,
		),
		compress1KB,
	))

	return mux
}
