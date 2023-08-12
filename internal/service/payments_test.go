package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	payment "github.com/grpc-buf/internal/gen/payment"
	"github.com/grpc-buf/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakePayment(t *testing.T) {
	mockDataStore := new(mocks.MockDataStore)
	paymentService := NewPaymentService(mockDataStore)

	req := &connect.Request[payment.PaymentRequest]{
		Msg: &payment.PaymentRequest{
			CardNo:         0,
			Card:           0,
			Name:           "",
			Address:        nil,
			Amount:         0,
			PaymentCreated: nil,
		},
	}

	expectedResponse := &connect.Response[payment.PaymentResponse]{
		Msg: &payment.PaymentResponse{Response: nil},
	}

	mockDataStore.On("MakePayment", mock.Anything, req).Return(expectedResponse, nil)

	resp, err := paymentService.MakePayment(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)

	mockDataStore.AssertExpectations(t)
}
