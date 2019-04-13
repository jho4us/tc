// Package repo provides implementations of the domain repositories.
package repo

import (
	"github.com/jho4us/tc/test"
	"github.com/jho4us/tc/testdb"
	"github.com/jho4us/tc/testdb/dbconf"
	"github.com/jho4us/tc/testdb/sql"
)

type testRepository struct {
	Accessor testdb.Accessor
}

func (r *testRepository) Store(t *test.Test) (test.TID, error) {
	if t == nil {
		return "", test.ErrInvalidArgument
	}
	if len(t.ID) == 0 {
		id, err := r.Accessor.InsertTest(testdb.TestRecord{
			Name:    t.Name,
			Content: t.Content,
		})
		if err != nil {
			return "", err
		}
		return test.TID(id), nil
	}
	err := r.Accessor.UpsertTest(testdb.TestRecord{
		ID:      string(t.ID),
		Name:    t.Name,
		Content: t.Content,
	})
	if err != nil {
		return "", err
	}
	return t.ID, nil
}

func (r *testRepository) Find(id test.TID) (*test.Test, error) {
	if len(id) == 0 {
		return nil, test.ErrInvalidArgument
	}
	rets, err := r.Accessor.GetTest(string(id))
	if err != nil {
		return nil, err
	}
	if len(rets) != 1 {
		return nil, test.ErrUnknown
	}
	t := &test.Test{
		ID:      test.TID(rets[0].ID),
		Name:    rets[0].Name,
		Content: rets[0].Content,
	}
	return t, nil
}

func (r *testRepository) FindAll() ([]*test.Test, error) {
	rets, err := r.Accessor.GetTests()
	if err != nil {
		t := make([]*test.Test, 0)
		return t, err
	}
	t := make([]*test.Test, 0, len(rets))
	for _, val := range rets {
		t = append(t, &test.Test{
			ID:      test.TID(val.ID),
			Name:    val.Name,
			Content: val.Content,
		})
	}
	return t, nil
}

func (r *testRepository) Delete(id test.TID) error {
	if len(id) == 0 {
		return test.ErrInvalidArgument
	}
	return r.Accessor.DeleteTest(string(id))
}

// NewTestRepository returns a new instance of a tests repository.
func NewTestRepository(path string) (test.Repository, error) {
	db, err := dbconf.DBFromConfig(path)
	if err != nil {
		return nil, err
	}
	accessor := sql.NewAccessor(db)
	return &testRepository{
		Accessor: accessor,
	}, nil
}