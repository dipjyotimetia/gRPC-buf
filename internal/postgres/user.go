package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"time"

	"connectrpc.com/connect/v2"
	"github.com/golang-jwt/jwt/v5"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// HashPassword takes a plain-text password and returns a hashed password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *Store) LoginUser(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	// Normalize and validate inputs
	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	password := strings.TrimSpace(req.GetPassword())
	if email == "" || password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, "email and password are required")
	}

	// Lookup user
	var storedPassword string
	var userID string
	err := s.db.QueryRow(ctx,
		"SELECT id, password FROM users WHERE email = $1",
		email).Scan(&userID, &storedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Dummy bcrypt compare to mitigate user enumeration timing
			_ = bcrypt.CompareHashAndPassword([]byte("$2a$10$7EqJtq98hPqEX7fNZaFWoO5B7N9gDCwQ4G6k8cQK6Y2Z2WfQ8Y/5e"), []byte(password))
			return nil, connect.NewError(connect.CodeUnauthenticated, "invalid credentials")
		}
		slog.Error("database error during login", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "internal server error")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, "invalid credentials")
	}

	// Build JWT
	now := time.Now().UTC()
	expirationTime := now.Add(15 * time.Minute)
	jti, err := randomJTI()
	if err != nil {
		slog.Error("error generating JWT ID", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "error generating authentication token")
	}
	issuer := strings.TrimSpace(s.sec.JWTIssuer)
	if issuer == "" {
		issuer = "grpc-buf"
	}
	aud := strings.TrimSpace(s.sec.JWTAudience)
	var auds jwt.ClaimStrings
	if aud != "" {
		auds = jwt.ClaimStrings{aud}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   userID,
		Audience:  auds,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		ID:        jti,
	})

	// Determine signing key from configuration
	signKey := strings.TrimSpace(s.sec.JWTSecret)
	if signKey == "" {
		slog.Error("JWT secret missing")
		return nil, connect.NewError(connect.CodeInternal, "authentication not configured")
	}

	tokenString, err := token.SignedString([]byte(signKey))
	if err != nil {
		slog.Error("error signing JWT token", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "error generating authentication token")
	}

	response := &userv1.LoginResponse{AccessToken: tokenString}
	return response, nil
}

func (s *Store) RegisterUser(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	// Normalize and validate inputs
	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	first := strings.TrimSpace(req.GetFirstName())
	last := strings.TrimSpace(req.GetLastName())
	pass := strings.TrimSpace(req.GetPassword())

	if email == "" || pass == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, "email and password are required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, "invalid email format")
	}
	if err := validatePassword(pass); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err.Error())
	}

	// Hash password
	hashedPassword, err := HashPassword(pass)
	if err != nil {
		slog.Error("error hashing password", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "error processing password")
	}

	// Insert user
	var userID string
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (email, password, first_name, last_name)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		email, hashedPassword, first, last).Scan(&userID)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerr.Code == "23505" { // unique_violation
			return nil, connect.NewError(connect.CodeAlreadyExists, "account with this email already exists")
		}
		slog.Error("error inserting user", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "error creating user account")
	}

	response := &userv1.RegisterResponse{
		Id:         userID,
		CreateTime: timestamppb.Now(),
	}
	return response, nil
}

// validatePassword enforces minimal complexity requirements.
func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	var lower, upper, digit, symbol bool
	for _, r := range p {
		switch {
		case 'a' <= r && r <= 'z':
			lower = true
		case 'A' <= r && r <= 'Z':
			upper = true
		case '0' <= r && r <= '9':
			digit = true
		default:
			symbol = true
		}
	}
	cnt := 0
	for _, b := range []bool{lower, upper, digit, symbol} {
		if b {
			cnt++
		}
	}
	if cnt < 3 {
		return errors.New("password must include a mix of letters, numbers, or symbols")
	}
	return nil
}

func randomJTI() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("randomJTI: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}
