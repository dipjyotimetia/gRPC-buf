package mongo

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	payment "github.com/grpc-buf/internal/gen/payment"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment struct {
	cardNO   int64            `bson:"cardNo"`
	cardType payment.CardType `bson:"card_type"`
	name     string           `bson:"name"`
	address  string           `bson:"address"`
	amount   float32          `bson:"amount"`
}

func (db *Store) MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error) {
	var data *Payment
	name := req.Msg.GetName()
	amount := req.Msg.GetAmount()
	if VerifyCard(req.Msg.CardNo) && req.Msg.Amount >= 0.0 {
		data = &Payment{
			cardNO:   req.Msg.CardNo,
			cardType: payment.CardType_DebitCard,
			name:     req.Msg.Name,
			address:  req.Msg.Address[0],
			amount:   req.Msg.Amount,
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
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID: %v", err))
	}

	log.Println("Got a payment with a name ", name, "amount", amount)
	log.Println(req.Header().Get("Some-Header"))
	response := connect.NewResponse(&payment.PaymentResponse{
		Response: &payment.PaymentResponse_Paid{Paid: true},
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}

func VerifyCard(cardNo int64) bool {
	return cardNo != 0
}
