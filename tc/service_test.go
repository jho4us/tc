package tc

import (
	"fmt"

	"github.com/jho4us/tc/repo"
	"github.com/jho4us/tc/test"
	"github.com/jho4us/tc/testdb/dbconf"

	"testing"
)

func TestService(t *testing.T) {
	rep, err := testCreateTestRepo(t)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	tcs := NewService(rep)
	_, err = tcs.CreateTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("unexistent")
	if err == nil {
		t.Fatal(err)
	}
	err = tcs.DeleteTest("")
	if err == nil {
		t.Fatal(err)
	}
	ta := tcs.Tests()
	if ta == nil {
		t.Errorf("cant query all tests")
	}
	testInsertTestAndGetTest(tcs, t, 0)
	testUpdateAndGetTest(tcs, t)
}

func testCreateTestRepo(t *testing.T) (test.Repository, error) {
	db, err := dbconf.DBFromConfig("./../repo/repo-tst-config.json")
	if err != nil {
		t.Fatal(err)
	}
	tr, err := repo.NewTestRepository(db)
	if err != nil {
		t.Fatal(err)
	}
	return tr, nil
}

func testInsertTestAndGetTest(s Service, t *testing.T, postfix int) test.ID {
	if postfix == 0 {
		testTruncateAll(s, t)
	}
	want := Test{
		Name:    fmt.Sprintf("fake name%d", postfix),
		Content: fmt.Sprintf("fake content%d", postfix),
	}

	id, err := s.CreateTest(want.Name)
	if err != nil {
		t.Fatal(err)
	}

	want.ID = string(id)
	err = s.PutTest(&want)

	if err != nil {
		t.Fatal(err)
	}

	got, err := s.LoadTest(id)
	if err != nil {
		t.Fatal(err)
	}

	// reflection comparison with zero time objects are not stable as it seems
	if want.Name != got.Name || want.Content != got.Content || want.ID != got.ID {
		t.Errorf("want TestRecord %+v, got %+v", want, got)
	}
	return id
}

func testTruncateAll(s Service, t *testing.T) {
	rets := s.Tests()
	for _, val := range rets {
		err := s.DeleteTest(test.ID(val.ID))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testUpdateAndGetTest(s Service, t *testing.T) {
	id := testInsertTestAndGetTest(s, t, 999)

	err := s.PutTest(nil)
	if err == nil {
		t.Errorf("null id put succeeded")
	}
	want := Test{
		ID:      "",
		Name:    fmt.Sprintf("fake name%d", 999),
		Content: fmt.Sprintf("fake content%d", 1000),
	}
	err = s.PutTest(&want)
	if err == nil {
		t.Errorf("null id put succeeded")
	}
	want.ID = string(id)
	err = s.PutTest(&want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.LoadTest(id)
	if err != nil {
		t.Fatal(err)
	}
	// reflection comparison with zero time objects are not stable as it seems
	if want.Name != got.Name || want.Content != got.Content || want.ID != got.ID {
		t.Errorf("want TestRecord %+v, got %+v", want, got)
	}
}
