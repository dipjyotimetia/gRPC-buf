package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"github.com/grpc-buf/internal/gen/proto/payment/paymentv1connect"

	"github.com/grpc-buf/internal/mongo"
	"github.com/grpc-buf/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestElizaServer(t *testing.T) {
	t.Parallel()
	var (
		db             = mongo.NewDatabaseConnection()
		paymentService = service.NewPaymentService(db)
	)
	mux := http.NewServeMux()
	mux.Handle(paymentv1connect.NewPaymentHandler(paymentService))
	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	connectClient := paymentv1connect.NewPaymentClient(
		server.Client(),
		server.URL,
	)
	grpcClient := paymentv1connect.NewPaymentClient(
		server.Client(),
		server.URL,
		connect.WithGRPC(),
	)
	clients := []paymentv1connect.PaymentClient{connectClient, grpcClient}

	t.Run("say", func(t *testing.T) {
		for _, client := range clients {
			result, err := client.MakePayment(context.Background(), connect.NewRequest(&paymentv1.PaymentRequest{
				CardNo:       123567887,
				Card:         2,
				Name:         "TestCard",
				AddressLines: []string{"efwefew"},
				Amount:       10,
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
