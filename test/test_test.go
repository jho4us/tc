package test

import (
	"github.com/pborman/uuid"

	"testing"
)

func TestCore(t *testing.T) {
	tr := New(TID("aaa"), "bbb", "ccc")
	if tr.ID != "aaa" || tr.Name != "bbb" || tr.Content != "ccc" {
		t.Errorf("structure data doesnt match")
	}
	tid := NextTestID()

	if uuid.Parse(string(tid)) == nil {
		t.Errorf("wrong uuid generated")

	}
}
