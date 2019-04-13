// Package testdb declare test entity ISP
package testdb

import (
	"time"
)

// TestRecord encodes a test and its metadata
// that will be recorded in a database.
type TestRecord struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Created    time.Time `db:"created_at"`
	ModifiedAt time.Time `db:"modified_at"`
	Content    string    `db:"content"`
}

// Accessor abstracts the CRUD of testdb objects from a DB.
type Accessor interface {
	InsertTest(tr TestRecord) (id string, err error)
	UpsertTest(tr TestRecord) error
	GetTest(id string) ([]TestRecord, error)
	GetTests() ([]TestRecord, error)
	DeleteTest(id string) error
}
