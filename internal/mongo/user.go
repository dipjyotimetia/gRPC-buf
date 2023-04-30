package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	userv1 "github.com/grpc-buf/internal/gen/registration"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	id        string `bson:"_id"`
	email     string `bson:"email"`
	password  string `bson:"password"`
	firstName string `bson:"first_name"`
	lastName  string `bson:"last_name"`
	createdAt string `bson:"createdAt"`
	updatedAt string `bson:"updatedAt"`
}

// hashPassword takes a plain-text password and returns a hashed password using bcrypt.
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (db *Store) LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	var result User
	filter := bson.D{{Key: "email", Value: req.Msg.GetEmail()}}
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "User not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(result.password), []byte(req.Msg.GetPassword())) != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Incorrect password")
	}

	expirationTime := time.Now().Add(15 * time.Minute) // expiration time of the token ->15 mins
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "grpc-buff",
		Subject:   "grpc-buff",
		Audience:  jwt.ClaimStrings{"grpc-buff"},
		ExpiresAt: &jwt.NumericDate{Time: expirationTime},
	})
	tokenString, err := token.SignedString("jwtKey")
	if err != nil {
		return nil, err
	}
	response := connect.NewResponse(&userv1.LoginResponse{
		AccessToken: tokenString,
	})
	return response, nil
}

func (db *Store) RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error) {
	var data *User
	hashedPassword, err := hashPassword(req.Msg.GetPassword())
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	data = &User{
		id:        uuid.New().String(),
		email:     req.Msg.GetEmail(),
		password:  hashedPassword,
		firstName: req.Msg.GetFirstName(),
		lastName:  req.Msg.GetLastName(),
		createdAt: timestamppb.Now().String(),
		updatedAt: timestamppb.Now().String(),
	}

	res, err := db.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID: %v", err))
	}

	response := connect.NewResponse(&userv1.RegisterResponse{
		Id:        id.Hex(),
		CreatedAt: nil,
	})
	response.Header().Set("Some-Other-Header", "hello!")
	return response, nil
}
