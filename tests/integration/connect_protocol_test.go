//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// This test uses the default Connect protocol (no WithGRPC/WithGRPCWeb option).
func TestConnectProtocol_MakePayment(t *testing.T) {
	client := paymentv1connect.NewPaymentServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	req := connect.NewRequest(&paymentv1.PaymentRequest{
		CardNo:       4242424242424242,
		Card:         paymentv1.CardType_CARD_TYPE_CREDIT,
		Name:         "Connect Test",
		AddressLines: []string{"123 Anywhere St"},
		Amount:       12.34,
		PaymentCreated: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	})
	if _, err := client.MakePayment(context.Background(), req); err != nil {
		t.Fatalf("connect protocol call failed: %v", err)
	}
}
