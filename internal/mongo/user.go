package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v4"
	userv1 "github.com/grpc-buf/internal/gen/registration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Email     string              `bson:"email"`
	Password  string              `bson:"password"`
	FirstName string              `bson:"firstName"`
	LastName  string              `bson:"lastName"`
	CreatedAt primitive.Timestamp `bson:"createdAt"`
	UpdatedAt primitive.Timestamp `bson:"updatedAt"`
}

type LoginCredentials struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// HashPassword takes a plain-text password and returns a hashed password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (db *Store) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	var result User
	email := req.Msg.GetEmail()
	filter := bson.D{{Key: "email", Value: email}}
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "User not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Msg.GetPassword())) != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Incorrect password")
	}

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
		return nil, err
	}
	response := connect.NewResponse(&userv1.LoginResponse{
		AccessToken: tokenString,
	})
	return response, nil
}

func (db *Store) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	hashedPassword, err := HashPassword(req.Msg.GetPassword())
	if err != nil {
		slog.Error("Error hashing password: %v", err)
	}

	data := User{
		ID:        primitive.NewObjectID(),
		Email:     req.Msg.GetEmail(),
		Password:  hashedPassword,
		FirstName: req.Msg.GetFirstName(),
		LastName:  req.Msg.GetLastName(),
		CreatedAt: primitive.Timestamp{T: uint32(time.Now().Unix())},
		UpdatedAt: primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	res, err := db.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID")
	}

	response := connect.NewResponse(&userv1.RegisterResponse{
		Id:        id.Hex(),
		CreatedAt: timestamppb.Now(),
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}
