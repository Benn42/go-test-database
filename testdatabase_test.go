package gotestdatabase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseInitialisation(t *testing.T) {
	tables := []string{"test", "tableone", "rainbow", "chicken_pie"}
	db := NewTestDatabase(tables)

	dbTables := db.ListTables()
	assert.Equal(t, len(dbTables), len(tables))

	for k, v := range dbTables {
		assert.Equal(t, tables[k], v)
	}
}

func TestDatabaseInsertOne(t *testing.T) {
	tables := []string{"test"}
	db := NewTestDatabase(tables)

	key := db.InsertOne("test", &DatabaseRecord{
		"test":  "one",
		"field": 3,
		"flag":  false,
	})

	result, err := db.FindOne("test", key)
	assert.Equal(t, err, nil)

	record := *result
	assert.Equal(t, record["test"], "one")
	assert.Equal(t, record["field"], 3)
	assert.Equal(t, record["flag"], false)
}

func TestDatabaseRemoveOne(t *testing.T) {
	tables := []string{"test"}
	db := NewTestDatabase(tables)

	key := db.InsertOne("test", &DatabaseRecord{
		"name": "remove",
		"age":  3,
	})

	_, err := db.FindOne("test", key)
	assert.Equal(t, err, nil)

	db.RemoveOne("test", key)
	_, rerr := db.FindOne("test", key)
	assert.NotNil(t, rerr)
	assert.Equal(t, rerr.Error(), "Record does not exist")
}
