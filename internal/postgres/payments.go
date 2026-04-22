package postgres

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MakePayment validates the payment request and records it as a debit-card
// payment. Low-level database errors are logged server-side and surfaced to
// the client as a generic Internal status to avoid leaking driver detail.
func (s *Store) MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	cardNo := req.Msg.GetCardNo()
	address := req.Msg.GetAddressLines()

	if !VerifyCard(cardNo) || amount < 0.0 || len(address) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "field validation failed")
	}

	_, err := s.db.Exec(ctx,
		`INSERT INTO payments (card_no, card_type, name, address, amount)
         VALUES ($1, $2, $3, $4, $5)`,
		cardNo, int(paymentv1.CardType_CARD_TYPE_DEBIT), name, address[0], amount)
	if err != nil {
		slog.Error("Error storing payment", "error", err)
		return nil, status.Error(codes.Internal, "failed to store payment")
	}

	slog.Info("Received payment", "name", name, "amount", amount)
	slog.Info("Header value", "value", req.Header().Get("Some-Header"))

	response := connect.NewResponse(&paymentv1.PaymentResponse{
		Status: paymentv1.PaymentStatus_PAYMENT_STATUS_PAID,
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

// VerifyCard performs a minimal sanity check on a card number. A zero-valued
// card number is treated as missing.
func VerifyCard(cardNo int64) bool {
	return cardNo != 0
}
