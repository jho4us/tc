package testingdb

import (
	"testing"
)

/* just to make tests happy */
func TestDb(t *testing.T) {
	_ = SQLiteDB("tstore_development.db")
}
