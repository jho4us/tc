// Package test contains the core of the domain model.
package test

import (
	"strings"

	"github.com/pborman/uuid"
)

// TID uniquely identifies a particular test.
type ID string

// Test is the central class in the domain model.
type Test struct {
	ID      ID
	Name    string
	Content string
}

// New creates a new test.
func New(id ID, name string, content string) *Test {
	return &Test{
		ID:      id,
		Name:    name,
		Content: content,
	}
}

// Repository provides access a test store.
type Repository interface {
	Store(t *Test) (ID, error)
	Find(id ID) (*Test, error)
	FindAll() ([]*Test, error)
	Delete(id ID) error
}

// NextTestID generates a new test ID.
func NextTestID() ID {
	return ID(strings.ToUpper(uuid.New()))
}
