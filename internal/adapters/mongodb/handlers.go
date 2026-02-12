package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
//	collections - list of collections names.
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

	names, err := m.db.ListCollectionNames(ctx, bson.M{})
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
	collection := m.db.Collection(collectionName)

	err := collection.Drop(context.Background())
	if err != nil {
		return fmt.Errorf("failed to drop collection: <%w>", err)
	}

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

	names, err = m.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %v", err)
	}

	return names, nil
}

// Send the document user on DB. Returns id added document and error.
//
// Params:
//
//	collectionName - name of collection
//	doc - document
func (m *mongoDB) SendDocumentUser(collectionName string, doc DocUser) (id interface{}, err error) {

	// Check
	if m.db == nil {
		return nil, ErrNilPtrDB
	}
	if m.connect == nil {
		return nil, ErrNilPtrConnect
	}
	if collectionName == "" {
		return nil, ErrEmptyCollectionsName
	}
	if doc.Age <= 0 && doc.Email == "" && doc.Name == "" {
		return nil, ErrEmptyDocument
	}
	if doc.Age <= 0 {
		return nil, ErrValueAge
	}

	// Сheck exists document
	collection := m.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var existingDoc DocUser
	filter := bson.M{"name": doc.Name, "age": doc.Age, "email": doc.Email}
	err = collection.FindOne(ctx, filter).Decode(&existingDoc)

	if err == nil {
		return nil, ErrDocumentExists
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("Function FindOne, return error: <%w>", err)
	}

	// Send
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("Function InsertOne, returned error: <%w>", err)
	}

	return result.InsertedID, nil
}

// Update the document user on DB. Returns id added document and error.
//
// Params:
//
//	collectionName - name of collection
//	name - name of user
//	doc - document
func (m *mongoDB) UpdateDocumentUserByName(collectionName, name string, doc DocUser) (err error) {

	// Check
	if m.db == nil {
		return ErrNilPtrDB
	}
	if m.connect == nil {
		return ErrNilPtrConnect
	}
	if collectionName == "" {
		return ErrEmptyCollectionsName
	}
	if name == "" {
		return ErrEmptyValueName
	}
	if doc.Age <= 0 && doc.Email == "" && doc.Name == "" {
		return ErrEmptyDocument
	}
	if doc.Age <= 0 {
		return ErrValueAge
	}

	// Logic
	collection := m.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"name": name}

	update := bson.M{"$set": doc}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Function UpdateOne, returned error: <%w>", err)
	}

	if result.MatchedCount == 0 {
		return ErrUpdateDocument
	}

	return nil
}

// Recieve document user by name. Returns document and error.
//
// Params:
//
//	collectionName - name of collection
//	name - name
func (m *mongoDB) RecvDocumentUserByName(collectionName string, name string) (doc DocUser, err error) {

	// Check
	if m.db == nil {
		return DocUser{}, ErrNilPtrDB
	}
	if m.connect == nil {
		return DocUser{}, ErrNilPtrConnect
	}
	if collectionName == "" {
		return DocUser{}, ErrEmptyCollectionsName
	}
	if name == "" {
		return DocUser{}, ErrEmptyValueName
	}

	// Logic
	collection := m.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"name": name}

	err = collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return DocUser{}, fmt.Errorf("Function FindOne return error: <%w>", err)
	}

	return doc, nil
}

// Delete document user by name. Returns document and error.
//
// Params:
//
//	collectionName - name of collection
//	name - name
func (m *mongoDB) DelDocumentUserByName(collectionName string, name string) (int64, error) {

	// Check
	if m.db == nil {
		return 0, ErrNilPtrDB
	}
	if m.connect == nil {
		return 0, ErrNilPtrConnect
	}
	if collectionName == "" {
		return 0, ErrEmptyCollectionsName
	}
	if name == "" {
		return 0, ErrEmptyValueName
	}

	// Logic
	collection := m.db.Collection(collectionName)

	filter := bson.M{"name": name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete document: <%w>", err)
	}

	return result.DeletedCount, nil
}

// Change collection for document. Return error.
//
// Params:
//
//	srcCollection - source collection
//	destCollection - destination collection
//	doc - document
func (m *mongoDB) MoveDocumentUser(srcCollection, destCollection string, doc DocUser) error {

	//
	// Check
	//

	if m.db == nil {
		return ErrNilPtrDB
	}
	if m.connect == nil {
		return ErrNilPtrConnect
	}
	if srcCollection == "" || destCollection == "" {
		return ErrEmptyCollectionsName
	}
	if doc.Name == "" {
		return ErrEmptyValueName
	}

	//
	// Logic
	//

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sourceCollection := m.db.Collection(srcCollection)
	destinationCollection := m.db.Collection(destCollection)
	filter := bson.M{"name": doc.Name}

	var result bson.M

	// Recieve document
	err := sourceCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("Document is not found: <%v>", err)
		}
		return fmt.Errorf("Fault recieve document: <%w>", err)
	}

	// Send document
	_, err = destinationCollection.InsertOne(ctx, result)
	if err != nil {
		return fmt.Errorf("Fault send document: <%w>", err)
	}

	// Delete document
	_, err = sourceCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Fault delete document: <%w>", err)
	}

	return nil
}
