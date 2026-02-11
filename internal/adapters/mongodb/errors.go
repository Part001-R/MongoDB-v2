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
	// Missing names of collections
	ErrEmptyCollectionsNames = errors.New("Missing names of collections")
	// Not correct DSN
	ErrNotCorrectDSN = errors.New("Not correct DSN")
)
