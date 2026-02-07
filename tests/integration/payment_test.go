//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPayment(t *testing.T) {
	client := paymentv1connect.NewPaymentServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)
	req := connect.NewRequest(&paymentv1.PaymentRequest{
		CardNo:       123567887,
		Card:         2,
		Name:         "TestCard",
		AddressLines: []string{"efwefew"},
		Amount:       10,
		PaymentCreated: &timestamppb.Timestamp{
			Seconds: int64(time.Now().Second()),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	})
	res, err := client.MakePayment(context.Background(), req)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("MakePayment failed: %v", err)
		return
	}
	fmt.Println(res.Msg)
}
