package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/grpc-buf/internal/gen/payment/paymentconnect"
	"golang.org/x/net/http2"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	"github.com/grpc-buf/internal/gen/payment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	client := paymentconnect.NewPaymentClient(
		&http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
					// If you're also using this client for non-h2c traffic, you may want to
					// delegate to tls.Dial if the network isn't TCP or the addr isn't in an
					// allowlist.
					return net.Dial(network, addr)
				},
			},
		},
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
	req.Header().Set("Some-Header", "hello from connect")
	res, err := client.MakePayment(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	log.Println(res.Header().Get("Some-Other-Header"))
}
