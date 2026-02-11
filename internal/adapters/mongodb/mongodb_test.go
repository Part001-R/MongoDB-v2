package mongodb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test New.
func TestNew(t *testing.T) {

	t.Run("Missing DSN", func(t *testing.T) {
		dsn := ""

		_, err := New(dsn)
		require.Equalf(t, ErrEmptyValueDSN, err, "Error is not equal")
	})

	t.Run("Correct", func(t *testing.T) {

		dsn := "mongodb://localhost:27017/myDatabase"

		db, err := New(dsn)
		require.NoErrorf(t, err, "Unexpected error New")
		require.NotNil(t, db, "Pointer db is nil")

		err = db.Close()
		require.NoErrorf(t, err, "Unexpected error Close")
	})
}
