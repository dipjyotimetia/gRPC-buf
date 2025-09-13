package postgres

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	constant "github.com/grpc-buf/internal/const"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (db *Store) MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	ctx, span := constant.Tracer.Start(ctx, "MakePayment")
	defer span.End()

	// Extract payment details from request
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	cardNo := req.Msg.GetCardNo()
	address := req.Msg.GetAddressLines()

	// Validate payment details
	if !VerifyCard(cardNo) || amount < 0.0 || len(address) == 0 {
		span.RecordError(status.Error(codes.FailedPrecondition, "field validation failed"))
		return nil, status.Errorf(codes.FailedPrecondition, "field validation failed")
	}

	// Insert payment into database
	_, err := db.db.Exec(ctx,
		`INSERT INTO payments (card_no, card_type, name, address, amount)
         VALUES ($1, $2, $3, $4, $5)`,
		cardNo, int(paymentv1.CardType_CARD_TYPE_DEBIT), name, address[0], amount)
	if err != nil {
		span.RecordError(err)
		slog.Error("Error storing payment", "error", err)
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	slog.Info("Received payment", "name", name, "amount", amount)
	slog.Info("Header value", "value", req.Header().Get("Some-Header"))

	response := connect.NewResponse(&paymentv1.PaymentResponse{
		Status: paymentv1.PaymentStatus_PAYMENT_STATUS_PAID,
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

func VerifyCard(cardNo int64) bool {
	return cardNo != 0
}
