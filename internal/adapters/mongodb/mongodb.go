package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Presentation
type mongoDB struct {
	db *mongo.Client
}

// Interface
type MongoDBI interface {
	// Close db connect
	Close() error
}

// Constructor.
func New(dsn string) (MongoDBI, error) {

	// Check
	if dsn == "" {
		return nil, ErrEmptyValueDSN
	}

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %v", err)
	}

	// Check connect
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %v", err)
	}

	// Instance
	return &mongoDB{
		db: client,
	}, nil
}
