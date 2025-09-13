package postgres

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "log/slog"
    "net/mail"
    "strings"
    "time"

    "connectrpc.com/connect"
    "github.com/golang-jwt/jwt/v5"
    userv1 "github.com/grpc-buf/internal/gen/proto/registration"
    "github.com/jackc/pgconn"
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
    // Normalize and validate inputs
    email := strings.ToLower(strings.TrimSpace(req.Msg.GetEmail()))
    password := strings.TrimSpace(req.Msg.GetPassword())
    if email == "" || password == "" {
        return nil, status.Error(codes.InvalidArgument, "email and password are required")
    }

    // Lookup user
    var storedPassword string
    var userID string
    err := db.db.QueryRow(ctx,
        "SELECT id, password FROM users WHERE email = $1",
        email).Scan(&userID, &storedPassword)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // Dummy bcrypt compare to mitigate user enumeration timing
            _ = bcrypt.CompareHashAndPassword([]byte("$2a$10$7EqJtq98hPqEX7fNZaFWoO5B7N9gDCwQ4G6k8cQK6Y2Z2WfQ8Y/5e"), []byte(password))
            return nil, status.Error(codes.Unauthenticated, "invalid credentials")
        }
        slog.Error("database error during login", "error", err)
        return nil, status.Error(codes.Internal, "internal server error")
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    }

    // Build JWT
    now := time.Now().UTC()
    expirationTime := now.Add(15 * time.Minute)
    jti := randomJTI()
    issuer := strings.TrimSpace(db.sec.JWTIssuer)
    if issuer == "" { issuer = "grpc-buf" }
    aud := strings.TrimSpace(db.sec.JWTAudience)
    var auds jwt.ClaimStrings
    if aud != "" { auds = jwt.ClaimStrings{aud} }
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
    signKey := strings.TrimSpace(db.sec.JWTSecret)
    if signKey == "" {
        slog.Error("JWT secret missing")
        return nil, status.Error(codes.Internal, "authentication not configured")
    }

    tokenString, err := token.SignedString([]byte(signKey))
    if err != nil {
        slog.Error("error signing JWT token", "error", err)
        return nil, status.Error(codes.Internal, "error generating authentication token")
    }

    response := connect.NewResponse(&userv1.LoginResponse{AccessToken: tokenString})
    return response, nil
}

func (db *Store) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
    // Normalize and validate inputs
    email := strings.ToLower(strings.TrimSpace(req.Msg.GetEmail()))
    first := strings.TrimSpace(req.Msg.GetFirstName())
    last := strings.TrimSpace(req.Msg.GetLastName())
    pass := strings.TrimSpace(req.Msg.GetPassword())

    if email == "" || pass == "" {
        return nil, status.Error(codes.InvalidArgument, "email and password are required")
    }
    if _, err := mail.ParseAddress(email); err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid email format")
    }
    if err := validatePassword(pass); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    // Hash password
    hashedPassword, err := HashPassword(pass)
    if err != nil {
        slog.Error("error hashing password", "error", err)
        return nil, status.Error(codes.Internal, "error processing password")
    }

    // Insert user
    var userID string
    err = db.db.QueryRow(ctx,
        `INSERT INTO users (email, password, first_name, last_name)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
        email, hashedPassword, first, last).Scan(&userID)
    if err != nil {
        var pgerr *pgconn.PgError
        if errors.As(err, &pgerr) && pgerr.Code == "23505" { // unique_violation
            return nil, status.Error(codes.AlreadyExists, "account with this email already exists")
        }
        slog.Error("error inserting user", "error", err)
        return nil, status.Error(codes.Internal, "error creating user account")
    }

    response := connect.NewResponse(&userv1.RegisterResponse{
        Id:        userID,
        CreatedAt: timestamppb.Now(),
    })
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

func randomJTI() string {
    var b [16]byte
    if _, err := rand.Read(b[:]); err != nil {
        return ""
    }
    return hex.EncodeToString(b[:])
}
