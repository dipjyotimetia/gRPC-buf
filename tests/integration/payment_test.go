//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPayment(t *testing.T) {
	client := paymentv1connect.NewPaymentServiceClient(connect.NewClient(connecthttp.NewTransport(
		http.DefaultClient,
		"http://localhost:8080",
		connecthttp.WithGRPC(),
	)))
	req := &paymentv1.PaymentRequest{
		CardToken:    "tok_integration_debit_1",
		Card:         paymentv1.CardType_CARD_TYPE_DEBIT,
		Name:         "TestCard",
		AddressLines: []string{"123 Test Street"},
		Amount:       &money.Money{CurrencyCode: "USD", Units: 10},
		PaymentCreated: &timestamppb.Timestamp{
			Seconds: int64(time.Now().Second()),
			Nanos:   int32(time.Now().Nanosecond()),
		},
	}
	res, err := client.MakePayment(context.Background(), req)
	require.NoError(t, err, "MakePayment failed")
	t.Logf("payment response: %+v", res)
}
