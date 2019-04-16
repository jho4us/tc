// Package sql implement Accessor interface to store data in SQL databases
package sql

import (
	"fmt"
	"strings"

	"github.com/jho4us/tc/test"
	"github.com/jho4us/tc/testdb"

	"github.com/jmoiron/sqlx"
	"github.com/kisielk/sqlstruct"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

// Match to sqlx
var (
	_ = func() struct{} {
		sqlstruct.TagName = "db"
		return struct{}{}
	}()
)

const (
	insertSQL = `
INSERT INTO tests (id, name, content)
	VALUES (:id, :name, :content);`

	selectSQL = `
SELECT %s FROM tests
	WHERE (id = ?);`

	selectAllSQL = `
SELECT %s FROM tests;`

	updateSQL = `
UPDATE tests
    SET name=:name, content=:content
    WHERE (id = :id);`

	deleteSQL = `
DELETE FROM tests
	WHERE (id = ?);`
)

// Accessor implements testdb.Accessor interface.
type Accessor struct {
	db *sqlx.DB
}

func wrapSQLError(err error) error {
	if err != nil {
		return errors.Wrap(err, "SQL execution error")
	}
	return nil
}

func (d *Accessor) checkDB() error {
	if d.db == nil {
		return errors.Errorf("unknown db object, please check SetDB method")
	}
	return nil
}

// NewAccessor returns a new Accessor.
func NewAccessor(db *sqlx.DB) *Accessor {
	return &Accessor{db: db}
}

// SetDB changes the underlying sql.DB object Accessor is manipulating.
func (d *Accessor) SetDB(db *sqlx.DB) {
	d.db = db
}

// InsertTest puts a testdb.TestRecord into db.
func (d *Accessor) InsertTest(tr testdb.TestRecord) (id string, err error) {
	err = d.checkDB()
	if err != nil {
		return "", err
	}
	id = tr.ID
	if len(id) == 0 {
		id = strings.ToUpper(uuid.New())
	}
	newName := strings.TrimSpace(tr.Name)
	if len(newName) == 0 {
		return "", test.ErrInvalidArgument
	}
	res, err := d.db.NamedExec(insertSQL, &testdb.TestRecord{
		ID:      id,
		Name:    newName,
		Content: tr.Content,
	})
	if err != nil {
		return "", wrapSQLError(err)
	}

	numRowsAffected, err := res.RowsAffected()

	if numRowsAffected == 0 {
		return "", errors.Errorf("failed to insert the test record")
	}

	if numRowsAffected != 1 {
		return "", wrapSQLError(errors.Errorf("%d rows are affected, should be 1 row", numRowsAffected))
	}

	return id, err
}

// GetTest gets a testdb.TestRecord indexed by id.
func (d *Accessor) GetTest(id string) (trs []testdb.TestRecord, err error) {
	err = d.checkDB()
	if err != nil {
		return nil, err
	}

	err = d.db.Select(&trs, fmt.Sprintf(d.db.Rebind(selectSQL), sqlstruct.Columns(testdb.TestRecord{})), id)
	if err != nil {
		return nil, wrapSQLError(err)
	}

	return trs, nil
}

// GetTests gets all test records from db.
func (d *Accessor) GetTests() (trs []testdb.TestRecord, err error) {
	err = d.checkDB()
	if err != nil {
		return nil, err
	}

	err = d.db.Select(&trs, fmt.Sprintf(d.db.Rebind(selectAllSQL), sqlstruct.Columns(testdb.TestRecord{})))
	if err != nil {
		return nil, wrapSQLError(err)
	}

	return trs, nil
}

// DeleteTest delete test record with a given id from db.
func (d *Accessor) DeleteTest(id string) error {
	err := d.checkDB()
	if err != nil {
		return err
	}

	result, err := d.db.Exec(deleteSQL, id)
	if err != nil {
		return wrapSQLError(err)
	}

	numRowsAffected, err := result.RowsAffected()

	if numRowsAffected == 0 {
		return errors.Errorf("failed to delete test: test %s not found", id)
	}

	if numRowsAffected != 1 {
		return wrapSQLError(fmt.Errorf("%d rows are affected, should be 1 row", numRowsAffected))
	}

	return err
}

// UpsertTest update a test record with a given id,
// or insert the record if it doesn't yet exist in the db
// Implementation note:
// We didn't implement 'upsert' with SQL statement and we lost race condition
// prevention provided by underlying DBMS.
// Reasoning:
// 1. it's diffcult to support multiple DBMS backends in the same time, the
// SQL syntax differs from one to another.
// 2. we don't need a strict simultaneous consistency
func (d *Accessor) UpsertTest(tr testdb.TestRecord) error {
	err := d.checkDB()
	if err != nil {
		return err
	}
	if len(tr.ID) == 0 {
		_, err = d.InsertTest(tr)
		return err
	}

	result, err := d.db.NamedExec(updateSQL, &testdb.TestRecord{
		ID:      tr.ID,
		Name:    tr.Name,
		Content: tr.Content,
	})

	if err != nil {
		return wrapSQLError(err)
	}

	numRowsAffected, err := result.RowsAffected()

	if numRowsAffected == 0 {
		_, err = d.InsertTest(tr)
		return err
	}

	if numRowsAffected != 1 {
		return wrapSQLError(fmt.Errorf("%d rows are affected, should be 1 row", numRowsAffected))
	}

	return err
}
