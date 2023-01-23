package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/grpc-buf/internal/gen/payment/paymentconnect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.

	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/grpc-buf/internal/gen/payment"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}

func run() error {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"payment.v1.MakePayment",
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
	)
	// The generated constructors return a path and a plain net/http
	// handler.
	log.Println("starting new server")
	mux.Handle(paymentconnect.NewPaymentHandler(&PaymentServer{}))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err) //nolint:gocritic
	}
	return nil
}

// PaymentServer implements the PetStoreService API.
type PaymentServer struct {
	paymentconnect.UnimplementedPaymentHandler
}

// MakePayment adds the pet associated with the given request into the PetStore.
func (s *PaymentServer) MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error) {
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	log.Println("Got a payment with a name ", name, "amount", amount)
	log.Println(req.Header().Get("Some-Header"))
	res := connect.NewResponse(&payment.PaymentResponse{
		Response: &payment.PaymentResponse_Paid{Paid: true},
	})
	res.Header().Set("Some-Other-Header", "hello!")

	return res, nil
}
