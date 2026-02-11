package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Close db connect. Return error.
func (d *mongoDB) Close() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := d.connect.Disconnect(ctx); err != nil {
		return fmt.Errorf("Function Disconnetc, return error <%w>", err)
	}

	return nil
}

// If DB is not exists - create by name. Return error.
//
// Params:
//
//	collections - list of collections name.
func (m *mongoDB) CheckCreateDB(collections []string) error {

	// Check
	if m.nameDB == "" {
		return ErrEmptyValueNameDB
	}
	if collections == nil {
		return ErrNilPtrCollections
	}
	if len(collections) == 0 {
		return ErrEmptyCollectionsNames
	}

	// Get collections names
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	database := m.connect.Database(m.nameDB)

	names, err := database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collection names: %v", err)
	}

	// Create collections
	for _, v := range collections {

		if IsExistsCollection(names, v) {
			continue
		}

		collection := m.connect.Database(m.nameDB).Collection(v)

		doc := bson.D{
			{Key: "name", Value: "initial"},
		}

		_, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			return fmt.Errorf("failed to create collection and insert initial document: %v", err)
		}
	}

	return nil
}

// Drop collection by name. Return error.
//
// Params:
//
//	collectionName - collection name.
func (m *mongoDB) DropCollection(collectionName string) error {
	// Check
	if collectionName == "" {
		return ErrEmptyValueName
	}
	if m.nameDB == "" {
		return ErrEmptyValueNameDB
	}

	// Logic
	collection := m.connect.Database(m.nameDB).Collection(collectionName)

	// Логируем название коллекции
	log.Printf("Attempting to drop collection: %s from database: %s", collectionName, m.nameDB)

	err := collection.Drop(context.Background())
	if err != nil {
		return fmt.Errorf("failed to drop collection: <%w>", err)
	}

	log.Printf("Successfully dropped collection: %s", collectionName)
	return nil
}

// Get names of collections. Returns names and error.
func (m *mongoDB) GetNamesCollections() (names []string, err error) {

	// Check
	if m.connect == nil {
		return nil, ErrNilPtrDB
	}
	if m.nameDB == "" {
		return nil, ErrEmptyValueNameDB
	}

	// Logic
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	database := m.connect.Database(m.nameDB)

	names, err = database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %v", err)
	}

	return names, nil
}
