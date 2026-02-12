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
		assert.Equalf(t, collections[1], rxAfter[0], "Names is not equals")
	})

}

// Test SendDocumentUser
func TestSendDocumentUser(t *testing.T) {

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

	defer func() {
		err := db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 0")

		err = db.DropCollection(collections[1])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 1")
	}()

	t.Run("Missing collection name", func(t *testing.T) {

		collection := ""
		doc := DocUser{
			Name:  "Aaa",
			Age:   20,
			Email: "AAA@mail.com",
		}

		_, err := db.SendDocumentUser(collection, doc)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Wrong age", func(t *testing.T) {

		collection := collections[0]
		doc := DocUser{
			Name:  "Aaa",
			Age:   -1,
			Email: "AAA@mail.com",
		}

		_, err := db.SendDocumentUser(collection, doc)
		require.Equalf(t, ErrValueAge, err, "Error is not equal")
	})

	t.Run("Missing document", func(t *testing.T) {

		collection := collections[0]
		doc := DocUser{}

		_, err := db.SendDocumentUser(collection, doc)
		require.Equalf(t, ErrEmptyDocument, err, "Error is not equal")
	})

	t.Run("Correct", func(t *testing.T) {

		collection := collections[0]
		doc := DocUser{
			Name:  "Aaa",
			Age:   30,
			Email: "AAA@mail.com",
		}

		_, err := db.SendDocumentUser(collection, doc)
		require.NoErrorf(t, err, "Unexpected error")
	})

	t.Run("Exists entry", func(t *testing.T) {

		collection := collections[0]
		doc := DocUser{
			Name:  "B",
			Age:   30,
			Email: "B@mail.com",
		}

		_, err := db.SendDocumentUser(collection, doc)
		require.NoErrorf(t, err, "Unexpected error")

		_, err = db.SendDocumentUser(collection, doc)
		require.Equalf(t, ErrDocumentExists, err, "Error is not equal")
	})

}

// Test UpdateDocumentUser
func TestUpdateDocumentUser(t *testing.T) {

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

	defer func() {
		err := db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 0")

		err = db.DropCollection(collections[1])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 1")
	}()

	t.Run("Missing collection name", func(t *testing.T) {

		collection := ""
		name := "Aaa"
		doc := DocUser{
			Name:  "Aaa",
			Age:   20,
			Email: "AAA@mail.com",
		}

		err := db.UpdateDocumentUserByName(collection, name, doc)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Missing name", func(t *testing.T) {

		collection := collections[0]
		name := ""
		doc := DocUser{
			Name:  "Aaa",
			Age:   20,
			Email: "AAA@mail.com",
		}

		err := db.UpdateDocumentUserByName(collection, name, doc)
		require.Equalf(t, ErrEmptyValueName, err, "Error is not equal")
	})

	t.Run("Missing document", func(t *testing.T) {

		collection := collections[0]
		name := "Aaa"
		doc := DocUser{}

		err := db.UpdateDocumentUserByName(collection, name, doc)
		require.Equalf(t, ErrEmptyDocument, err, "Error is not equal")
	})

	t.Run("Not correct age", func(t *testing.T) {

		collection := collections[0]
		name := "Aaa"
		doc := DocUser{
			Name:  "Aaa",
			Age:   0,
			Email: "Bbb@mail.com",
		}

		err := db.UpdateDocumentUserByName(collection, name, doc)
		require.Equalf(t, ErrValueAge, err, "Error is not equal")
	})

	t.Run("Correct", func(t *testing.T) {

		// Send
		collection := collections[0]
		doc := DocUser{
			Name:  "Aaa",
			Age:   30,
			Email: "Bbb@mail.com",
		}
		db.SendDocumentUser(collection, doc)
		require.NoErrorf(t, err, "Unexpected error send")

		// Update
		name := "Aaa"
		doc2 := DocUser{
			Name:  "Aaa",
			Age:   33,
			Email: "Bbb@mail.com",
		}
		err := db.UpdateDocumentUserByName(collection, name, doc2)
		require.NoErrorf(t, err, "Unexpected error update")

		// Recieve
		rxDoc, err := db.RecvDocumentUserByName(collection, name)
		require.NoErrorf(t, err, "Unexpected error recieve")

		// Check
		assert.Equalf(t, doc2.Name, rxDoc.Name, "Value name is not equal")
		assert.Equalf(t, doc2.Age, rxDoc.Age, "Value age is not equal")
		assert.Equalf(t, doc2.Email, rxDoc.Email, "Value email is not equal")
	})
}

// Test RecvDocumentUserByName
func TestRecvDocumentUserByName(t *testing.T) {

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

	defer func() {
		err := db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 0")

		err = db.DropCollection(collections[1])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 1")
	}()

	t.Run("Missing collection name", func(t *testing.T) {

		collection := ""
		name := "AAA"

		_, err := db.RecvDocumentUserByName(collection, name)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Missing name", func(t *testing.T) {

		collection := collections[0]
		name := ""

		_, err := db.RecvDocumentUserByName(collection, name)
		require.Equalf(t, ErrEmptyValueName, err, "Error is not equal")
	})

	t.Run("Missing document", func(t *testing.T) {

		collection := collections[0]
		name := "B"

		_, err := db.RecvDocumentUserByName(collection, name)
		require.Equalf(t, "Function FindOne return error: <mongo: no documents in result>", err.Error(), "Error is not equal")

	})

	t.Run("Correct", func(t *testing.T) {

		// Add
		collection := collections[0]
		doc := DocUser{
			Name:  "Aaa",
			Age:   20,
			Email: "AAA@mail.com",
		}
		_, err := db.SendDocumentUser(collection, doc)
		require.NoErrorf(t, err, "Unexpected error send")

		// Recieve
		rxDoc, err := db.RecvDocumentUserByName(collection, doc.Name)
		require.NoErrorf(t, err, "Unexpected error recieve")

		// Check
		assert.Equalf(t, doc.Name, rxDoc.Name, "Name is not equal")
		assert.Equalf(t, doc.Age, rxDoc.Age, "Age is not equal")
		assert.Equalf(t, doc.Email, rxDoc.Email, "Email is not equal")
	})
}

// Test DelDocumentUserByName
func TestDelDocumentUserByName(t *testing.T) {

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

	defer func() {
		err := db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 0")

		err = db.DropCollection(collections[1])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 1")
	}()

	t.Run("Missing collection name", func(t *testing.T) {

		collection := ""
		name := "AAA"

		_, err := db.DelDocumentUserByName(collection, name)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Missing name", func(t *testing.T) {

		collection := collections[0]
		name := ""

		_, err := db.DelDocumentUserByName(collection, name)
		require.Equalf(t, ErrEmptyValueName, err, "Error is not equal")
	})

	t.Run("Not exists entry", func(t *testing.T) {

		collection := collections[0]
		name := "CCC"

		cntRel, err := db.DelDocumentUserByName(collection, name)
		require.NoErrorf(t, err, "Unexpected error")
		assert.Equalf(t, int64(0), cntRel, "Value is not equal")

	})

	t.Run("Correct", func(t *testing.T) {

		collection := collections[0]
		doc := DocUser{
			Name:  "Ccc",
			Age:   30,
			Email: "Ccc@mail.com",
		}

		// Send
		_, err := db.SendDocumentUser(collection, doc)
		require.NoErrorf(t, err, "Unexpected error send")

		// Delete
		cntDel, err := db.DelDocumentUserByName(collection, doc.Name)
		require.NoErrorf(t, err, "Unexpected error delete")
		assert.Equalf(t, int64(1), cntDel, "Value is not equal")

	})
}

// Test MoveDocumentUser
func TestMoveDocumentUser(t *testing.T) {

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

	defer func() {
		err := db.DropCollection(collections[0])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 0")

		err = db.DropCollection(collections[1])
		require.NoErrorf(t, err, "Unexpected error CheckCreateDB 1")
	}()

	t.Run("Missing srcCollection name", func(t *testing.T) {

		srcCollection := ""
		destCollection := collections[1]
		doc := DocUser{
			Name:  "A",
			Age:   30,
			Email: "A@mail.mail",
		}

		err := db.MoveDocumentUser(srcCollection, destCollection, doc)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Missing destCollection name", func(t *testing.T) {

		srcCollection := collections[0]
		destCollection := ""
		doc := DocUser{
			Name:  "A",
			Age:   30,
			Email: "A@mail.mail",
		}

		err := db.MoveDocumentUser(srcCollection, destCollection, doc)
		require.Equalf(t, ErrEmptyCollectionsName, err, "Error is not equal")
	})

	t.Run("Missing name", func(t *testing.T) {

		srcCollection := collections[0]
		destCollection := collections[1]
		doc := DocUser{
			Name:  "",
			Age:   30,
			Email: "A@mail.mail",
		}

		err := db.MoveDocumentUser(srcCollection, destCollection, doc)
		require.Equalf(t, ErrEmptyValueName, err, "Error is not equal")
	})

	t.Run("Correct", func(t *testing.T) {

		srcCollection := collections[0]
		destCollection := collections[1]
		doc := DocUser{
			Name:  "A",
			Age:   30,
			Email: "A@mail.mail",
		}

		// Create
		_, err := db.SendDocumentUser(srcCollection, doc)
		require.NoErrorf(t, err, "Unexpected error send")

		// Check exists
		_, err = db.RecvDocumentUserByName(srcCollection, doc.Name)
		require.NoErrorf(t, err, "Document was not found after creation")

		// Relocating
		err = db.MoveDocumentUser(srcCollection, destCollection, doc)
		require.NoErrorf(t, err, "Unexpected error move")

		// Receive
		rxDoc, err := db.RecvDocumentUserByName(destCollection, doc.Name)
		require.NoErrorf(t, err, "Unexpected error receive")

		// Check
		assert.Equalf(t, doc.Name, rxDoc.Name, "Name is not equal")
		assert.Equalf(t, doc.Age, rxDoc.Age, "Age is not equal")
		assert.Equalf(t, doc.Email, rxDoc.Email, "Email is not equal")

		_, err = db.RecvDocumentUserByName(srcCollection, doc.Name)
		require.Error(t, err, "Document should be deleted from source collection")
	})

}
