//go:build integration
// +build integration

package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"

	"github.com/grpc-buf/internal/postgres"
	"github.com/grpc-buf/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestElizaServer(t *testing.T) {
	t.Parallel()

	// Skip if DATABASE_URL is not set (e.g., when testing against docker-compose API)
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping test: DATABASE_URL not set")
	}

	db, err := postgres.NewDatabaseConnection(context.Background())
	require.NoError(t, err)
	defer db.Close()
	paymentService := service.NewPaymentService(db)
	mux := http.NewServeMux()
	mux.Handle(paymentv1connect.NewPaymentServiceHandler(paymentService))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := paymentv1connect.NewPaymentServiceClient(
		server.Client(),
		server.URL,
	)
	grpcClient := paymentv1connect.NewPaymentServiceClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)
	clients := []paymentv1connect.PaymentServiceClient{connectClient, grpcClient}

	t.Run("say", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.MakePayment(context.Background(), connect.NewRequest(&paymentv1.PaymentRequest{
				CardToken:    "tok_test_debit_1",
				Card:         paymentv1.CardType_CARD_TYPE_DEBIT,
				Name:         "TestCard",
				AddressLines: []string{"123 Test Street"},
				Amount:       &money.Money{CurrencyCode: "USD", Units: 10},
				PaymentCreated: &timestamppb.Timestamp{
					Seconds: int64(time.Now().Second()),
					Nanos:   int32(time.Now().Nanosecond()),
				},
			}))
			require.Nil(t, err)
			assert.True(t, len(result.Msg.String()) > 0)
		}
	})
}
