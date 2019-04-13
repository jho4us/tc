package repo

import (
	"fmt"

	"github.com/jho4us/tc/test"

	"testing"
)

func TestRepo(t *testing.T) {
	_, err := NewTestRepository("nonexistent")
	if err == nil {
		t.Fatal(err)
	}
	tr, err := NewTestRepository("./repo-tst-config.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tr.Store(nil)
	if err == nil {
		t.Fatal(err)
	}
	_, err = tr.Find("")
	if err == nil {
		t.Fatal(err)
	}
	err = tr.Delete("")
	if err == nil {
		t.Fatal(err)
	}

	testInsertTestAndGetTest(tr, t, 0)
	testUpdateAndGetTest(tr, t)
	testGetTests(tr, t)
}

func testInsertTestAndGetTest(tr test.Repository, t *testing.T, postfix int) test.TID {
	if postfix == 0 {
		testTruncateAll(tr, t)
	}
	want := test.Test{
		Name:    fmt.Sprintf("fake name%d", postfix),
		Content: fmt.Sprintf("fake content%d", postfix),
	}

	id, err := tr.Store(&want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := tr.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Errorf("test %s not found", id)
	}

	// reflection comparison with zero time objects are not stable as it seems
	if want.Name != got.Name || want.Content != got.Content || got.ID != id {
		t.Errorf("want TestRecord %+v, got %+v", want, got)
	}
	return id
}

func testUpdateAndGetTest(tr test.Repository, t *testing.T) {
	id := testInsertTestAndGetTest(tr, t, 999)

	want := test.Test{
		ID:      id,
		Name:    fmt.Sprintf("fake name%d", 999),
		Content: fmt.Sprintf("fake content%d", 1000),
	}
	newID, err := tr.Store(&want)
	if err != nil {
		t.Fatal(err)
	}

	if newID != id {
		t.Errorf("id doesnt match: want %s got %s", id, newID)

	}

	got, err := tr.Find(newID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Errorf("test %s not found", id)
	}

	// reflection comparison with zero time objects are not stable as it seems
	if want.Name != got.Name || want.Content != got.Content || got.ID != id {
		t.Errorf("want TestRecord %+v, got %+v", want, got)
	}
}

func testTruncateAll(tr test.Repository, t *testing.T) {
	rets, err := tr.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, val := range rets {
		err := tr.Delete(val.ID)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testGetTests(tr test.Repository, t *testing.T) {
	testTruncateAll(tr, t)

	for i := 1; i <= 10; i++ {
		testInsertTestAndGetTest(tr, t, i)
	}
	rets, err := tr.FindAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(rets) != 10 {
		t.Fatal("should return 10 records.")
	}
}
