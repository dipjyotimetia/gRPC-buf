package postgres

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect/v2"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"google.golang.org/genproto/googleapis/type/money"
)

// MakePayment validates the payment request and records it as a debit-card
// payment. Low-level database errors are logged server-side and surfaced to
// the client as a generic Internal status to avoid leaking driver detail.
func (s *Store) MakePayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	name := req.GetName()
	amount := req.GetAmount()
	cardToken := strings.TrimSpace(req.GetCardToken())
	address := req.GetAddressLines()

	if !VerifyCardToken(cardToken) || !positiveMoney(amount) || len(address) == 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, "field validation failed")
	}

	amountCents := amount.GetUnits()*100 + int64(amount.GetNanos())/nanosPerCent
	currency := amount.GetCurrencyCode()
	_, err := s.db.Exec(ctx,
		`INSERT INTO payments (card_token, card_type, name, address, amount_cents, currency_code)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		cardToken, int(paymentv1.CardType_CARD_TYPE_DEBIT), name, address[0], amountCents, currency)
	if err != nil {
		slog.Error("Error storing payment", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "failed to store payment")
	}

	slog.Info("Received payment", "name", name, "amount_cents", amountCents, "currency", currency)
	// CallInfo is absent when the store is invoked outside a Connect handler
	// (e.g. via the MCP adapters), so header access is best-effort.
	info, hasInfo := connect.CallInfoForServerContext(ctx)
	if hasInfo {
		slog.Info("Header value", "value", info.RequestHeader().Get("Some-Header"))
	}

	response := &paymentv1.PaymentResponse{
		Status: paymentv1.PaymentStatus_PAYMENT_STATUS_PAID,
	}
	if hasInfo {
		info.ResponseHeader().Set("Some-Other-Header", "hello!")
	}
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
