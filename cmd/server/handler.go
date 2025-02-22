package server

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"github.com/grpc-buf/internal/gen/proto/registration/userv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func setupHandler() *http.ServeMux {
	setupOtel()

	mux := http.NewServeMux()

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		log.Fatal(err)
	}

	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(paymentv1connect.NewPaymentHandler(
		paymentService,
		compress1KB,
		connect.WithInterceptors(otelInterceptor),
	))

	mux.Handle(userv1connect.NewUserServiceHandler(
		userService,
		compress1KB,
		connect.WithInterceptors(otelInterceptor),
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

func setupOtel() {
	otel.SetTracerProvider(trace.NewTracerProvider())
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
