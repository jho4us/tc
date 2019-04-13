// Package test contains the core of the domain model.
package test

import (
	"strings"

	"github.com/pborman/uuid"
)

// TID uniquely identifies a particular test.
type TID string

// Test is the central class in the domain model.
type Test struct {
	ID      TID
	Name    string
	Content string
}

// New creates a new test.
func New(id TID, name string, content string) *Test {
	return &Test{
		ID:      id,
		Name:    name,
		Content: content,
	}
}

// Repository provides access a test store.
type Repository interface {
	Store(t *Test) (TID, error)
	Find(id TID) (*Test, error)
	FindAll() ([]*Test, error)
	Delete(id TID) error
}

// NextTestID generates a new test ID.
func NextTestID() TID {
	return TID(strings.ToUpper(uuid.New()))
}
