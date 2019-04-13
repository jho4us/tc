package sql

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/jho4us/tc/testdb"
	"github.com/jho4us/tc/testdb/testingdb"

	"github.com/jmoiron/sqlx"
)

const (
	sqliteDBFile = "../testingdb/tstore_development.db"
)

func TestNoDB(t *testing.T) {
	dba := &Accessor{}
	_, err := dba.GetTest("a89a220f-c5ab-41ae-854a-b22f89dfdfdf")
	if err == nil {
		t.Fatal("should return error")
	}
	dba.SetDB(nil)
	utcTime := time.Now().UTC()
	want := testdb.TestRecord{
		Name:       fmt.Sprintf("fake name%d", 0),
		Content:    fmt.Sprintf("fake content%d", 0),
		Created:    utcTime,
		ModifiedAt: utcTime,
	}

	_, err = dba.InsertTest(want)
	if err == nil {
		t.Fatal("should return error")
	}

}

type TestAccessor struct {
	Accessor testdb.Accessor
	DB       *sqlx.DB
}

func (ta *TestAccessor) Truncate() {
	testingdb.Truncate(ta.DB)
}

func TestSQLite(t *testing.T) {
	db := testingdb.SQLiteDB(sqliteDBFile)
	ta := TestAccessor{
		Accessor: NewAccessor(db),
		DB:       db,
	}
	testEverything(ta, t)
}

// roughlySameTime decides if t1 and t2 are close enough.
func roughlySameTime(t1, t2 time.Time) bool {
	// return true if the difference is smaller than sec.
	return math.Abs(float64(t1.Sub(t2))) < float64(time.Second)
}

func testEverything(ta TestAccessor, t *testing.T) {
	testMisc(ta, t)
	testInsertTestAndGetTest(ta, t, 0)
	testUpsertTestAndGetTest(ta, t)
	testDeleteTest(ta, t)
	testGetTests(ta, t)
}

func testMisc(ta TestAccessor, t *testing.T) {
	if wrapSQLError(nil) != nil {
		t.Fatal("should not wrap nil erorr.")
	}
	err := errors.New("unknown test")
	if wrapSQLError(err) == err {
		t.Fatal("Error not wrapped.")
	}
}

func testInsertTestAndGetTest(ta TestAccessor, t *testing.T, postfix int) {
	if postfix == 0 {
		ta.Truncate()
	}

	utcTime := time.Now().UTC()
	want := testdb.TestRecord{
		Name:       fmt.Sprintf("fake name%d", postfix),
		Content:    fmt.Sprintf("fake content%d", postfix),
		Created:    utcTime,
		ModifiedAt: utcTime,
	}

	id, err := ta.Accessor.InsertTest(want)
	if err != nil {
		t.Fatal(err)
	}

	rets, err := ta.Accessor.GetTest(id)
	if err != nil {
		t.Fatal(err)
	}

	if len(rets) != 1 {
		t.Fatal("should only return one record.")
	}

	got := rets[0]

	// reflection comparison with zero time objects are not stable as it seems
	if want.Name != got.Name || want.Content != got.Content ||
		!roughlySameTime(got.Created, utcTime) || !roughlySameTime(got.ModifiedAt, utcTime) {
		t.Errorf("want TestRecord %+v, got %+v", want, got)
	}
}

func testUpsertTestAndGetTest(ta TestAccessor, t *testing.T) {
	ta.Truncate()

	utcTime := time.Now().UTC()
	want := testdb.TestRecord{
		Name:       "fake name",
		Content:    "fake content",
		Created:    utcTime,
		ModifiedAt: utcTime,
	}
	id, err := ta.Accessor.InsertTest(want)
	if err != nil {
		t.Fatal(err)
	}
	want.ID = id
	want.Name = "good name"
	want.Content = "good content"

	if err := ta.Accessor.UpsertTest(want); err != nil {
		t.Fatal(err)
	}

	rets, err := ta.Accessor.GetTest(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(rets) != 1 {
		t.Fatal("should return exactly one record")
	}

	got := rets[0]

	if want.Name != got.Name || want.Content != got.Content ||
		!roughlySameTime(want.Created, got.Created) ||
		want.ModifiedAt == got.ModifiedAt {
		t.Errorf("want OCSP %+v, got %+v", want, got)
	}
	want.ID = ""
	if err := ta.Accessor.UpsertTest(want); err != nil {
		t.Fatal(err)
	}
	want.ID = "nonexistent"
	if err := ta.Accessor.UpsertTest(want); err != nil {
		t.Fatal(err)
	}

}

func testDeleteTest(ta TestAccessor, t *testing.T) {
	ta.Truncate()

	utcTime := time.Now().UTC()
	want := testdb.TestRecord{
		Name:       "fake name",
		Content:    "fake content",
		Created:    utcTime,
		ModifiedAt: utcTime,
	}

	id, err := ta.Accessor.InsertTest(want)
	if err != nil {
		t.Fatal(err)
	}

	err = ta.Accessor.DeleteTest(id)
	if err != nil {
		t.Fatal(err)
	}

	trs, err := ta.Accessor.GetTest(id)
	if trs != nil {
		t.Fatal(err)
	}

	err = ta.Accessor.DeleteTest("nonexistent")
	if err == nil {
		t.Fatal(err)
	}

}

func testGetTests(ta TestAccessor, t *testing.T) {
	ta.Truncate()

	for i := 1; i <= 10; i++ {
		testInsertTestAndGetTest(ta, t, i)
	}
	rets, err := ta.Accessor.GetTests()
	if err != nil {
		t.Fatal(err)
	}

	if len(rets) != 10 {
		t.Fatal("should return 10 records.")
	}
}
