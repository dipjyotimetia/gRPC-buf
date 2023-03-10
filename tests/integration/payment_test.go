package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	paymentconnect "github.com/grpc-buf/internal/gen/payment/paymentv1connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	payment "github.com/grpc-buf/internal/gen/payment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPayment(t *testing.T) {
	setupOtel()
	client := paymentconnect.NewPaymentClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithInterceptors(otelconnect.NewInterceptor()),
		connect.WithGRPC(),
	)
	req := connect.NewRequest(&payment.PaymentRequest{
		CardNo:  123567887,
		Card:    2,
		Name:    "TestCard",
		Address: []string{"efwefew"},
		Amount:  10,
		PaymentCreated: &timestamppb.Timestamp{
			Seconds: int64(time.Now().Second()),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	})
	res, err := client.MakePayment(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.Msg)
	fmt.Println(res.Header().Get("Some-Other-Header"))
}

func setupOtel() {
	// Exporting to different platforms can be configured here
	otel.SetTracerProvider(trace.NewTracerProvider())
	global.SetMeterProvider(metric.NewMeterProvider())
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
