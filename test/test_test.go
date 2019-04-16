package test

import (
	"github.com/pborman/uuid"

	"testing"
)

func TestNew(t *testing.T) {
	tr := New(ID("aaa"), "bbb", "ccc")
	if tr.ID != "aaa" || tr.Name != "bbb" || tr.Content != "ccc" {
		t.Errorf("structure data doesnt match")
	}
	tid := NextTestID()

	if uuid.Parse(string(tid)) == nil {
		t.Errorf("wrong uuid generated")

	}
}
