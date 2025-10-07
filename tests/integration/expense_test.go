//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/grpc-buf/internal/gen/proto/expense/expensev1connect"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestExpenseCRUD(t *testing.T) {
	client := expensev1connect.NewExpenseServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	// Create
	createReq := connect.NewRequest(&expensev1.CreateExpenseRequest{
		Expense: &expensev1.Expense{
			UserId:      "00000000-0000-0000-0000-000000000001",
			Amount:      &money.Money{CurrencyCode: "USD", Units: 12},
			Category:    "test",
			Description: "integration",
		},
	})
	createRes, err := client.CreateExpense(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreateExpense failed: %v", err)
	}
	id := createRes.Msg.GetId()
	if id == "" {
		t.Fatalf("expected id from CreateExpense")
	}

	// Get
	getRes, err := client.GetExpense(context.Background(), connect.NewRequest(&expensev1.GetExpenseRequest{Id: id}))
	if err != nil {
		t.Fatalf("GetExpense failed: %v", err)
	}
	if getRes.Msg.GetId() != id {
		t.Fatalf("unexpected id from GetExpense")
	}

	// List
	listRes, err := client.ListExpenses(context.Background(), connect.NewRequest(&expensev1.ListExpensesRequest{UserId: getRes.Msg.GetUserId()}))
	if err != nil {
		t.Fatalf("ListExpenses failed: %v", err)
	}
	if len(listRes.Msg.GetExpenses()) == 0 {
		t.Fatalf("expected expenses in list")
	}

	// Update
	updReq := connect.NewRequest(&expensev1.UpdateExpenseRequest{
		Expense:    &expensev1.Expense{Id: id, Description: "updated"},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
	})
	_, err = client.UpdateExpense(context.Background(), updReq)
	if err != nil {
		t.Fatalf("UpdateExpense failed: %v", err)
	}

	// Delete
	_, err = client.DeleteExpense(context.Background(), connect.NewRequest(&expensev1.DeleteExpenseRequest{Id: id}))
	if err != nil {
		t.Fatalf("DeleteExpense failed: %v", err)
	}
}
