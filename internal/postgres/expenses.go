package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"connectrpc.com/connect"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	"github.com/jackc/pgx/v5"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const nanosPerCent = 10_000_000

func moneyToCents(m *money.Money) int64 {
	if m == nil {
		return 0
	}
	return m.GetUnits()*100 + int64(m.GetNanos())/nanosPerCent
}

func centsToMoney(cents int64, currency string) *money.Money {
	units := cents / 100
	rem := int32(cents % 100)
	return &money.Money{
		CurrencyCode: currency,
		Units:        units,
		Nanos:        rem * nanosPerCent,
	}
}

// CreateExpense inserts a new expense row and returns it with its generated id
// and timestamps.
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

	var (
		id                     string
		createTime, updateTime time.Time
	)
	err := s.db.QueryRow(ctx,
		`INSERT INTO expenses (user_id, amount_cents, currency_code, category, description)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
		exp.GetUserId(), moneyToCents(exp.GetAmount()), exp.GetAmount().GetCurrencyCode(), exp.GetCategory(), exp.GetDescription(),
	).Scan(&id, &createTime, &updateTime)
	if err != nil {
		slog.Error("create expense query failed", "error", err, "user_id", exp.GetUserId())
		return nil, status.Error(codes.Internal, "failed to create expense")
	}

	return connect.NewResponse(&expensev1.Expense{
		Id:          id,
		UserId:      exp.GetUserId(),
		Amount:      exp.GetAmount(),
		Category:    exp.GetCategory(),
		Description: exp.GetDescription(),
		CreateTime:  timestamppb.New(createTime),
		UpdateTime:  timestamppb.New(updateTime),
	}), nil
}

// GetExpense fetches a single expense by id. Returns codes.NotFound when the
// row does not exist.
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "expense not found")
		}
		slog.Error("get expense query failed", "error", err, "id", id)
		return nil, status.Error(codes.Internal, "failed to get expense")
	}

	return connect.NewResponse(&expensev1.Expense{
		Id:          id,
		UserId:      userID,
		Amount:      centsToMoney(amountCents, currency),
		Category:    category,
		Description: description,
		CreateTime:  timestamppb.New(createdAt),
		UpdateTime:  timestamppb.New(updatedAt),
	}), nil
}

// ListExpenses returns a page of expenses ordered by creation time (newest
// first), optionally filtered by user_id. Pagination uses opaque "o:<offset>"
// tokens.
func (s *Store) ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error) {
	userID := strings.TrimSpace(req.Msg.GetUserId())
	pageSize := req.Msg.GetPageSize()
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 50
	}

	offset := 0
	if req.Msg.GetPageToken() != "" {
		if n, err := fmt.Sscanf(req.Msg.GetPageToken(), "o:%d", &offset); n != 1 || err != nil {
			offset = 0
		}
	}

	var (
		rows pgx.Rows
		err  error
	)
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
		slog.Error("list expenses query failed", "error", err, "user_id", userID)
		return nil, status.Error(codes.Internal, "failed to list expenses")
	}
	defer rows.Close()

	resp := &expensev1.ListExpensesResponse{}
	for rows.Next() {
		var (
			id, uid, currency, category, description string
			amountCents                              int64
			createdAt, updatedAt                     time.Time
		)
		if err := rows.Scan(&id, &uid, &amountCents, &currency, &category, &description, &createdAt, &updatedAt); err != nil {
			slog.Error("list expenses scan failed", "error", err)
			return nil, status.Error(codes.Internal, "failed to list expenses")
		}
		resp.Expenses = append(resp.Expenses, &expensev1.Expense{
			Id:          id,
			UserId:      uid,
			Amount:      centsToMoney(amountCents, currency),
			Category:    category,
			Description: description,
			CreateTime:  timestamppb.New(createdAt),
			UpdateTime:  timestamppb.New(updatedAt),
		})
	}
	if err := rows.Err(); err != nil {
		slog.Error("list expenses iteration failed", "error", err)
		return nil, status.Error(codes.Internal, "failed to list expenses")
	}
	if len(resp.Expenses) == int(pageSize) {
		resp.NextPageToken = fmt.Sprintf("o:%d", offset+int(pageSize))
	}
	return connect.NewResponse(resp), nil
}

// UpdateExpense applies a field-mask update to an expense row. Supported
// mask paths: category, description, amount.
func (s *Store) UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error) {
	exp := req.Msg.GetExpense()
	if exp == nil || strings.TrimSpace(exp.GetId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "expense.id is required")
	}

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
		args = append(args, moneyToCents(exp.Amount), exp.Amount.GetCurrencyCode())
		idx += 2
	}
	if len(set) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask has no supported fields")
	}
	args = append(args, exp.GetId())

	_, err := s.db.Exec(ctx, fmt.Sprintf("UPDATE expenses SET %s, updated_at=NOW() WHERE id=$%d", strings.Join(set, ","), idx), args...)
	if err != nil {
		slog.Error("update expense query failed", "error", err, "id", exp.GetId())
		return nil, status.Error(codes.Internal, "failed to update expense")
	}
	return s.GetExpense(ctx, connect.NewRequest(&expensev1.GetExpenseRequest{Id: exp.GetId()}))
}

// DeleteExpense removes an expense by id. Returns codes.NotFound when no row
// matched. Returns an empty response on success (AIP-135).
func (s *Store) DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[emptypb.Empty], error) {
	id := strings.TrimSpace(req.Msg.GetId())
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	res, err := s.db.Exec(ctx, `DELETE FROM expenses WHERE id=$1`, id)
	if err != nil {
		slog.Error("delete expense query failed", "error", err, "id", id)
		return nil, status.Error(codes.Internal, "failed to delete expense")
	}
	if res.RowsAffected() == 0 {
		return nil, status.Error(codes.NotFound, "expense not found")
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}
