package mongo

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	constant "github.com/grpc-buf/internal/const"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment struct {
	CardNo   int64              `bson:"cardNo"`
	CardType paymentv1.CardType `bson:"card_type"`
	Name     string             `bson:"name"`
	Address  string             `bson:"address"`
	Amount   float32            `bson:"amount"`
}

func (db *Store) MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error) {
	ctx, span := constant.Tracer.Start(ctx, "MakePayment")
	defer span.End()

	var data *Payment
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	cardNo := req.Msg.GetCardNo()
	address := req.Msg.GetAddressLines()
	if VerifyCard(cardNo) && amount >= 0.0 {
		data = &Payment{
			CardNo:   cardNo,
			CardType: paymentv1.CardType_DebitCard,
			Name:     name,
			Address:  address[0],
			Amount:   amount,
		}
	} else {
		span.RecordError(fmt.Errorf("field validation failed"))
		return nil, status.Errorf(codes.FailedPrecondition, "field validation failed")
	}

	res, err := db.InsertOne(ctx, data)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		span.RecordError(fmt.Errorf("Cannot convert to OID"))
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID")
	}

	slog.Info("Received payment with name:", name, "and amount:", amount)
	slog.Info(req.Header().Get("Some-Header"))
	response := connect.NewResponse(&paymentv1.PaymentResponse{
		Status: paymentv1.PaymentStatus_PAID,
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

func VerifyCard(cardNo int64) bool {
	return cardNo != 0
}
