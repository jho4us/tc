package dbconf

import (
	"testing"

	_ "github.com/mattn/go-sqlite3" // import just to initialize SQLite testingdb
)

func TestLoadFile(t *testing.T) {
	config, err := loadConfigurationFile("testdata/db-config.json")
	if err != nil || config == nil {
		t.Fatal("Failed to load test db-config file ", err)
	}

	config, err = loadConfigurationFile("nonexistent")
	if err == nil || config != nil {
		t.Fatal("Expected failure loading nonexistent configuration file")
	}

	config, err = loadConfigurationFile("")
	if err == nil || config != nil {
		t.Fatal("Expected failure loading null configuration file")
	}
}

func TestDBFromConfig(t *testing.T) {
	db, err := DBFromConfig("testdata/db-config.json")
	if err != nil || db == nil {
		t.Fatal("Failed to open db from test db-config file")
	}

	db, err = DBFromConfig("testdata/bad-db-config.json")
	if err == nil || db != nil {
		t.Fatal("Expected failure opening invalid db")
	}

	db, err = DBFromConfig("testdata/empty-db-config.json")
	if err == nil || db != nil {
		t.Fatal("Expected failure opening invalid db")
	}

	db, err = DBFromConfig("testdata/non-json-config.json")
	if err == nil || db != nil {
		t.Fatal("Expected failure opening invalid db")
	}

	db, err = DBFromConfig("testdata/unreachable-db-config.json")
	if err == nil || db != nil {
		t.Fatal("Expected failure opening unreachable db")
	}
}
