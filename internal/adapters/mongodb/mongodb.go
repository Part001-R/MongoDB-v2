package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Presentation
type mongoDB struct {
	connect *mongo.Client
	nameDB  string
}

// Interface
type MongoDBI interface {
	// Close db connect.
	Close() error
	// Check-create DB.
	CheckCreateDB(collections []string) error
	// Drop collection by name.
	DropCollection(collectionName string) error
	// Get names of collections.
	GetNamesCollections() (names []string, err error)
}

// Constructor.
func New(dsn string) (MongoDBI, error) {

	// Check
	if dsn == "" {
		return nil, ErrEmptyValueDSN
	}
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		return nil, ErrNotCorrectDSN
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
	nameDB := parts[len(parts)-1]

	return &mongoDB{
		connect: client,
		nameDB:  nameDB,
	}, nil
}
