package postgres

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MakePayment validates the payment request and records it as a debit-card
// payment. Low-level database errors are logged server-side and surfaced to
// the client as a generic Internal status to avoid leaking driver detail.
func (s *Store) MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	cardToken := strings.TrimSpace(req.Msg.GetCardToken())
	address := req.Msg.GetAddressLines()

	if !VerifyCardToken(cardToken) || !positiveMoney(amount) || len(address) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "field validation failed")
	}

	amountCents := amount.GetUnits()*100 + int64(amount.GetNanos())/nanosPerCent
	currency := amount.GetCurrencyCode()
	_, err := s.db.Exec(ctx,
		`INSERT INTO payments (card_token, card_type, name, address, amount_cents, currency_code)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		cardToken, int(paymentv1.CardType_CARD_TYPE_DEBIT), name, address[0], amountCents, currency)
	if err != nil {
		slog.Error("Error storing payment", "error", err)
		return nil, status.Error(codes.Internal, "failed to store payment")
	}

	slog.Info("Received payment", "name", name, "amount_cents", amountCents, "currency", currency)
	slog.Info("Header value", "value", req.Header().Get("Some-Header"))

	response := connect.NewResponse(&paymentv1.PaymentResponse{
		Status: paymentv1.PaymentStatus_PAYMENT_STATUS_PAID,
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

// VerifyCardToken performs a minimal sanity check on a card token. An empty
// token is treated as missing.
func VerifyCardToken(token string) bool {
	return token != ""
}

func positiveMoney(m *money.Money) bool {
	if m == nil {
		return false
	}
	if m.GetUnits() < 0 || m.GetNanos() < 0 {
		return false
	}
	return m.GetUnits() > 0 || m.GetNanos() > 0
}
