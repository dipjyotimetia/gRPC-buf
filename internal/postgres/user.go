package postgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	constant "github.com/grpc-buf/internal/const"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (db *Store) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	ctx, span := constant.Tracer.Start(ctx, "LoginUser")
	defer span.End()

	// Get user credentials from request
	email := req.Msg.GetEmail()
	password := req.Msg.GetPassword()

	// Query for the user
	var storedPassword string
	var userID string
	err := db.db.QueryRow(ctx,
		"SELECT id, password FROM users WHERE email = $1",
		email).Scan(&userID, &storedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unauthenticated, "User not found")
		}
		span.RecordError(err)
		slog.Error("Database error during login", "error", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	// Compare password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Incorrect password")
	}

	// Generate JWT token
	expirationTime := time.Now().Add(15 * time.Minute) // Expiration time of the token: 15 minutes
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "grpc-buff",
		Subject:   "grpc-buff",
		Audience:  jwt.ClaimStrings{"grpc-buff"},
		ExpiresAt: &jwt.NumericDate{Time: expirationTime},
		ID:        email,
	})

	tokenString, err := token.SignedString([]byte("jwtKey"))
	if err != nil {
		span.RecordError(err)
		slog.Error("Error signing JWT token", "error", err)
		return nil, status.Errorf(codes.Internal, "Error generating authentication token")
	}

	response := connect.NewResponse(&userv1.LoginResponse{
		AccessToken: tokenString,
	})
	return response, nil
}

func (db *Store) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	ctx, span := constant.Tracer.Start(ctx, "RegisterUser")
	defer span.End()

	// Hash the password
	hashedPassword, err := HashPassword(req.Msg.GetPassword())
	if err != nil {
		span.RecordError(err)
		slog.Error("Error hashing password", "error", err)
		return nil, status.Errorf(codes.Internal, "Error processing password")
	}

	// Insert user into database
	var userID string
	err = db.db.QueryRow(ctx,
		`INSERT INTO users (email, password, first_name, last_name)
         VALUES ($1, $2, $3, $4)
         RETURNING id, created_at`,
		req.Msg.GetEmail(), hashedPassword, req.Msg.GetFirstName(), req.Msg.GetLastName()).Scan(&userID)
	if err != nil {
		span.RecordError(err)
		slog.Error("Error inserting user", "error", err)
		return nil, status.Errorf(codes.Internal, "Error creating user account")
	}

	response := connect.NewResponse(&userv1.RegisterResponse{
		Id:        userID,
		CreatedAt: timestamppb.Now(),
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}
