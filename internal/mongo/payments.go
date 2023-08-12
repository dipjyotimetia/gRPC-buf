package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	payment "github.com/grpc-buf/internal/gen/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment struct {
	CardNo   int64            `bson:"cardNo"`
	CardType payment.CardType `bson:"card_type"`
	Name     string           `bson:"name"`
	Address  string           `bson:"address"`
	Amount   float32          `bson:"amount"`
}

func (db *Store) MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error) {
	var data *Payment
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	cardNo := req.Msg.GetCardNo()
	address := req.Msg.GetAddress()
	if VerifyCard(cardNo) && amount >= 0.0 {
		data = &Payment{
			CardNo:   cardNo,
			CardType: payment.CardType_DebitCard,
			Name:     name,
			Address:  address[0],
			Amount:   amount,
		}
	} else {
		return nil, status.Errorf(codes.FailedPrecondition, "field validation failed")
	}

	res, err := db.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID")
	}

	slog.Info("Received payment with name:", name, "and amount:", amount)
	slog.Info(req.Header().Get("Some-Header"))
	response := connect.NewResponse(&payment.PaymentResponse{
		Response: &payment.PaymentResponse_Paid{Paid: true},
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

func VerifyCard(cardNo int64) bool {
	return cardNo != 0
}
