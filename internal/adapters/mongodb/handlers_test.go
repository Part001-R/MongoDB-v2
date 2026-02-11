package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test CheckCreateDB
func TestCheckCreateDB(t *testing.T) {

	dsn := "mongodb://localhost:27017/myDatabase"

	db, err := New(dsn)
	require.NoErrorf(t, err, "Unexpected error New")
	require.NotNil(t, db, "Pointer db is nil")

	defer func() {
		err = db.Close()
		assert.NoErrorf(t, err, "Unexpected error Close")
	}()

	t.Run("Missing pointer collections name", func(t *testing.T) {

		err := db.CheckCreateDB(nil)
		require.Equalf(t, ErrNilPtrCollections, err, "Error is not exists")
	})

	t.Run("Missing names", func(t *testing.T) {

		collections := []string{}

		err := db.CheckCreateDB(collections)
		require.Equalf(t, ErrEmptyCollectionsNames, err, "Error is not exists")
	})

	t.Run("Correct", func(t *testing.T) {

		collections := []string{"info-1", "info-2"}

		err := db.CheckCreateDB(collections)
		require.NoErrorf(t, err, "Unexpected error")
	})

}

// Test DropCollection
func TestDropCollection(t *testing.T) {

	dsn := "mongodb://localhost:27017/myDatabase"

	db, err := New(dsn)
	require.NoErrorf(t, err, "Unexpected error New")
	require.NotNil(t, db, "Pointer db is nil")

	defer func() {
		err = db.Close()
		assert.NoErrorf(t, err, "Unexpected error Close")
	}()

	collections := []string{"info-1", "info-2"}

	err = db.CheckCreateDB(collections)
	require.NoErrorf(t, err, "Unexpected error CheckCreateDB")

	t.Run("Missing name", func(t *testing.T) {

		collection := ""

		err := db.DropCollection(collection)
		require.Equalf(t, ErrEmptyValueName, err, "Error is not exists")
	})

	t.Run("Correct", func(t *testing.T) {

		collection := collections[0]

		err := db.DropCollection(collection)
		require.NoErrorf(t, err, "Unexpected error DropCollection")
	})

}

// Test GetNamesCollections
func TestGetNamesCollections(t *testing.T) {

	dsn := "mongodb://localhost:27017/myDatabase"

	db, err := New(dsn)
	require.NoErrorf(t, err, "Unexpected error New")
	require.NotNil(t, db, "Pointer db is nil")

	defer func() {
		err = db.Close()
		assert.NoErrorf(t, err, "Unexpected error Close")
	}()

	collections := []string{"info-1", "info-2"}

	err = db.CheckCreateDB(collections)
	require.NoErrorf(t, err, "Unexpected error CheckCreateDB")

	t.Run("Correct", func(t *testing.T) {

		rxBefore, err := db.GetNamesCollections()
		require.NoErrorf(t, err, "Unexpected error GetNamesCollections before")

		err = db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error DropCollection")

		rxAfter, err := db.GetNamesCollections()
		require.NoErrorf(t, err, "Unexpected error GetNamesCollections after")

		assert.NotEqualf(t, 0, len(rxBefore), "Empty list before drop")
		assert.NotEqualf(t, 0, len(rxAfter), "Empty list after drop")
		assert.NotEqualf(t, len(rxBefore), len(rxAfter), "Values is equals")
	})

}
