package mongo

import (
	"context"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	payment "github.com/grpc-buf/internal/gen/payment"
	userv1 "github.com/grpc-buf/internal/gen/registration"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ENV = "ENVIRONMENT"

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type DataStore interface {
	MakePayment(ctx context.Context, req *connect.Request[payment.PaymentRequest]) (*connect.Response[payment.PaymentResponse], error)
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
}

// Store database session
type Store struct {
	*mongo.Collection
}

func NewDatabaseConnection() DataStore {
	env := os.Getenv(ENV)
	var err error
	var client *mongo.Client

	if env == "dev" {
		log.Info("Connecting to MongoDB local")
		client, err = mongo.NewClient(options.Client().
			SetAuth(options.Credential{
				Username: "admin",
				Password: "admin",
			}).ApplyURI("mongodb://localhost:27017").
			SetConnectTimeout(connectTimeout).
			SetMaxConnIdleTime(maxConnIdleTime).
			SetMinPoolSize(minPoolSize).
			SetMaxPoolSize(maxPoolSize))
	} else {
		log.Info("Connecting to MongoDB Atlas")
		client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	}

	if err != nil {
		log.Fatalf("MongoDB client creation failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	return &Store{
		Collection: client.Database("buf").Collection("grpc"),
	}
}
