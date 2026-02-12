package mongodb

import "errors"

var (
	// Empty value DSN
	ErrEmptyValueDSN = errors.New("Empty value DSN")
	// Empty value name
	ErrEmptyValueName = errors.New("Empty value name")
	// Empty value name DB
	ErrEmptyValueNameDB = errors.New("Empty value name DB")
	// Nil pointer collections
	ErrNilPtrCollections = errors.New("Nil pointer collections")
	// Nil pointer DB
	ErrNilPtrDB = errors.New("Nil pointer DB")
	// Nil pointer connect
	ErrNilPtrConnect = errors.New("Nil pointer connect")
	// Empty collections names
	ErrEmptyCollectionsNames = errors.New("Empty collections names")
	// Empty collection name
	ErrEmptyCollectionsName = errors.New("Empty collection name")
	// Empty document
	ErrEmptyDocument = errors.New("Empty document")
	// Not correct DSN
	ErrNotCorrectDSN = errors.New("Not correct DSN")
	// Document exists
	ErrDocumentExists = errors.New("Document exists")
	// Error value age
	ErrValueAge = errors.New("Error value age")
	// Error update document
	ErrUpdateDocument = errors.New("Error update document")
)
