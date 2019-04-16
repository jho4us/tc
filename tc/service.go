// Package tc provides the model of test constructor.
package tc

import (
	"strings"

	"github.com/jho4us/tc/test"
)

// Service is the interface that provides test contructor methods.
type Service interface {
	// CreateTest registers a new empty test in the system
	CreateTest(name string) (test.ID, error)

	// LoadTest returns a read model of a test.
	LoadTest(id test.ID) (Test, error)

	// PutTest update test contents
	PutTest(t *Test) error

	// DeleteTest delete test from a system.
	DeleteTest(id test.ID) error

	// Tests returns a list of all tests in a system.
	Tests() []Test
}

type service struct {
	tests test.Repository
}

func (s *service) CreateTest(name string) (test.ID, error) {
	newName := strings.TrimSpace(name)
	if len(newName) == 0 {
		return "", test.ErrInvalidArgument
	}

	id := test.NextTestID()
	t := test.New(id, newName, "")
	id, err := s.tests.Store(t)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) LoadTest(id test.ID) (Test, error) {
	if id == "" {
		return Test{}, test.ErrInvalidArgument
	}

	t, err := s.tests.Find(test.ID(strings.ToUpper(string(id))))
	if err != nil {
		return Test{}, err
	}

	return Test{ID: string(t.ID), Name: t.Name, Content: t.Content}, nil
}

func (s *service) PutTest(t *Test) error {
	if t == nil || t.ID == "" {
		return test.ErrInvalidArgument
	}
	tu := test.New(test.ID(t.ID), t.Name, t.Content)
	_, err := s.tests.Store(tu)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteTest(id test.ID) error {
	if id == "" {
		return test.ErrInvalidArgument
	}

	return s.tests.Delete(test.ID(strings.ToUpper(string(id))))
}

func (s *service) Tests() []Test {
	result := make([]Test, 0)
	rs, err := s.tests.FindAll()
	if err == nil {
		for _, t := range rs {
			result = append(result, Test{
				ID:      string(t.ID),
				Name:    t.Name,
				Content: t.Content,
			})
		}
	}
	return result
}

// NewService creates a test constructor service with necessary dependencies.
func NewService(ts test.Repository) Service {
	return &service{
		tests: ts,
	}
}

// Test is a read model for tc views.
type Test struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content"`
}
