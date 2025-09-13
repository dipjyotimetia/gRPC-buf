package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	money "github.com/grpc-buf/internal/gen/proto/google/type"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	exp := req.Msg.GetExpense()
	if exp == nil {
		return nil, status.Error(codes.InvalidArgument, "expense is required")
	}
	if strings.TrimSpace(exp.GetUserId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if exp.Amount == nil || exp.Amount.CurrencyCode == "" {
		return nil, status.Error(codes.InvalidArgument, "amount currency_code is required")
	}

	var id string
	var createTime time.Time
	var updateTime time.Time

	units := exp.GetAmount().GetUnits()
	nanos := exp.GetAmount().GetNanos()
	amountCents := units*100 + int64(nanos)/10_000_000 // 1e7 nanos per cent

	err := s.db.QueryRow(ctx,
		`INSERT INTO expenses (user_id, amount_cents, currency_code, category, description)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
		exp.GetUserId(), amountCents, exp.GetAmount().GetCurrencyCode(), exp.GetCategory(), exp.GetDescription(),
	).Scan(&id, &createTime, &updateTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create expense: %v", err)
	}

	resp := &expensev1.Expense{
		Id:          id,
		UserId:      exp.GetUserId(),
		Amount:      exp.GetAmount(),
		Category:    exp.GetCategory(),
		Description: exp.GetDescription(),
		CreateTime:  timestamppb.New(createTime),
		UpdateTime:  timestamppb.New(updateTime),
	}
	return connect.NewResponse(resp), nil
}

func (s *Store) GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	id := strings.TrimSpace(req.Msg.GetId())
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var (
		userID, currency, category, description string
		amountCents                             int64
		createdAt, updatedAt                    time.Time
	)
	err := s.db.QueryRow(ctx,
		`SELECT user_id, amount_cents, currency_code, category, description, created_at, updated_at
         FROM expenses WHERE id=$1`, id,
	).Scan(&userID, &amountCents, &currency, &category, &description, &createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "expense not found: %v", err)
	}

	units := amountCents / 100
	cents := amountCents % 100
	resp := &expensev1.Expense{
		Id:          id,
		UserId:      userID,
		Amount:      &money.Money{CurrencyCode: currency, Units: units, Nanos: int32(cents) * 10_000_000},
		Category:    category,
		Description: description,
		CreateTime:  timestamppb.New(createdAt),
		UpdateTime:  timestamppb.New(updatedAt),
	}
	return connect.NewResponse(resp), nil
}

func (s *Store) ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error) {
	userID := strings.TrimSpace(req.Msg.GetUserId())
	pageSize := req.Msg.GetPageSize()
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 50
	}

	// Simple pagination using offset encoded in page_token
	offset := 0
	if req.Msg.GetPageToken() != "" {
		if n, err := fmt.Sscanf(req.Msg.GetPageToken(), "o:%d", &offset); n != 1 || err != nil {
			offset = 0
		}
	}

	var rows pgx.Rows
	var err error
	if userID == "" {
		rows, err = s.db.Query(ctx,
			`SELECT id, user_id, amount_cents, currency_code, category, description, created_at, updated_at
             FROM expenses ORDER BY created_at DESC LIMIT $1 OFFSET $2`, pageSize, offset,
		)
	} else {
		rows, err = s.db.Query(ctx,
			`SELECT id, user_id, amount_cents, currency_code, category, description, created_at, updated_at
             FROM expenses WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, pageSize, offset,
		)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list query failed: %v", err)
	}
	defer rows.Close()

	resp := &expensev1.ListExpensesResponse{}
	for {
		var (
			id, uid, currency, category, description string
			amountCents                              int64
			createdAt, updatedAt                     time.Time
		)
		if !rows.Next() {
			break
		}
		if err := rows.Scan(&id, &uid, &amountCents, &currency, &category, &description, &createdAt, &updatedAt); err != nil {
			return nil, status.Errorf(codes.Internal, "scan failed: %v", err)
		}
		units := amountCents / 100
		cents := amountCents % 100
		resp.Expenses = append(resp.Expenses, &expensev1.Expense{
			Id:          id,
			UserId:      uid,
			Amount:      &money.Money{CurrencyCode: currency, Units: units, Nanos: int32(cents) * 10_000_000},
			Category:    category,
			Description: description,
			CreateTime:  timestamppb.New(createdAt),
			UpdateTime:  timestamppb.New(updatedAt),
		})
	}
	if len(resp.Expenses) == int(pageSize) {
		resp.NextPageToken = fmt.Sprintf("o:%d", offset+int(pageSize))
	}
	return connect.NewResponse(resp), nil
}

func (s *Store) UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	exp := req.Msg.GetExpense()
	if exp == nil || strings.TrimSpace(exp.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "expense.id is required")
	}

	// Build dynamic update set based on update_mask paths
	paths := map[string]bool{}
	if mask := req.Msg.GetUpdateMask(); mask != nil {
		for _, p := range mask.Paths {
			paths[p] = true
		}
	}

	set := []string{}
	args := []any{}
	idx := 1
	if paths["category"] {
		set = append(set, fmt.Sprintf("category=$%d", idx))
		args = append(args, exp.GetCategory())
		idx++
	}
	if paths["description"] {
		set = append(set, fmt.Sprintf("description=$%d", idx))
		args = append(args, exp.GetDescription())
		idx++
	}
	if paths["amount"] && exp.Amount != nil {
		set = append(set, fmt.Sprintf("amount_cents=$%d,currency_code=$%d", idx, idx+1))
		args = append(args, exp.Amount.GetUnits()*100+int64(exp.Amount.GetNanos())/10_000_000, exp.Amount.GetCurrencyCode())
		idx += 2
	}
	if len(set) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask has no supported fields")
	}
	args = append(args, exp.GetId())

	_, err := s.db.Exec(ctx, fmt.Sprintf("UPDATE expenses SET %s, updated_at=NOW() WHERE id=$%d", strings.Join(set, ","), idx), args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	// Return updated row
	return s.GetExpense(ctx, connect.NewRequest(&expensev1.GetExpenseRequest{Id: exp.GetId()}))
}

func (s *Store) DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[timestamppb.Timestamp], error) {
	id := strings.TrimSpace(req.Msg.GetId())
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	res, err := s.db.Exec(ctx, `DELETE FROM expenses WHERE id=$1`, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	n := res.RowsAffected()
	if n == 0 {
		return nil, status.Error(codes.NotFound, "expense not found")
	}
	return connect.NewResponse(timestamppb.Now()), nil
}
