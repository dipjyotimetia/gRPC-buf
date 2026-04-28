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
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// This test uses the default Connect protocol (no WithGRPC/WithGRPCWeb option).
func TestConnectProtocol_MakePayment(t *testing.T) {
	client := paymentv1connect.NewPaymentServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	req := connect.NewRequest(&paymentv1.PaymentRequest{
		CardToken:    "tok_connect_credit_1",
		Card:         paymentv1.CardType_CARD_TYPE_CREDIT,
		Name:         "Connect Test",
		AddressLines: []string{"123 Anywhere St"},
		Amount:       &money.Money{CurrencyCode: "USD", Units: 12, Nanos: 340_000_000},
		PaymentCreated: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	})
	if _, err := client.MakePayment(context.Background(), req); err != nil {
		t.Fatalf("connect protocol call failed: %v", err)
	}
}
