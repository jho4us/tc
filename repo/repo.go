// Package repo provides implementations of the domain repositories.
package repo

import (
	"github.com/jmoiron/sqlx"

	"github.com/jho4us/tc/test"
	"github.com/jho4us/tc/testdb"
	"github.com/jho4us/tc/testdb/sql"
)

type testRepository struct {
	accessor testdb.Accessor
}

func (r *testRepository) Store(t *test.Test) (test.ID, error) {
	if t == nil {
		return "", test.ErrInvalidArgument
	}
	if len(t.ID) == 0 {
		id, err := r.accessor.InsertTest(testdb.TestRecord{
			Name:    t.Name,
			Content: t.Content,
		})
		if err != nil {
			return "", err
		}
		return test.ID(id), nil
	}
	err := r.accessor.UpsertTest(testdb.TestRecord{
		ID:      string(t.ID),
		Name:    t.Name,
		Content: t.Content,
	})
	if err != nil {
		return "", err
	}
	return t.ID, nil
}

func (r *testRepository) Find(id test.ID) (*test.Test, error) {
	if len(id) == 0 {
		return nil, test.ErrInvalidArgument
	}
	rets, err := r.accessor.GetTest(string(id))
	if err != nil {
		return nil, err
	}
	if len(rets) != 1 {
		return nil, test.ErrUnknown
	}
	t := &test.Test{
		ID:      test.ID(rets[0].ID),
		Name:    rets[0].Name,
		Content: rets[0].Content,
	}
	return t, nil
}

func (r *testRepository) FindAll() ([]*test.Test, error) {
	rets, err := r.accessor.GetTests()
	if err != nil {
		t := make([]*test.Test, 0)
		return t, err
	}
	t := make([]*test.Test, 0, len(rets))
	for _, val := range rets {
		t = append(t, &test.Test{
			ID:      test.ID(val.ID),
			Name:    val.Name,
			Content: val.Content,
		})
	}
	return t, nil
}

func (r *testRepository) Delete(id test.ID) error {
	if len(id) == 0 {
		return test.ErrInvalidArgument
	}
	return r.accessor.DeleteTest(string(id))
}

// NewTestRepository returns a new instance of a tests repository.
func NewTestRepository(db *sqlx.DB) (test.Repository, error) {
	if db == nil {
		return nil, test.ErrInvalidArgument
	}
	accessor := sql.NewAccessor(db)
	return &testRepository{
		accessor: accessor,
	}, nil
}
