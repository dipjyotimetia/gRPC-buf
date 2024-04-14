package main

import (
	"context"
	"fmt"
	"log/slog"

	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	paymentconnect "github.com/grpc-buf/internal/gen/payment/paymentv1connect"

	payment "github.com/grpc-buf/internal/gen/payment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPayment(t *testing.T) {
	client := paymentconnect.NewPaymentClient(
		http.DefaultClient,
		"http://localhost:8080",
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
		slog.Error(err.Error())
	}
	fmt.Println(res.Msg)
	fmt.Println(res.Header().Get("Some-Other-Header"))
}
