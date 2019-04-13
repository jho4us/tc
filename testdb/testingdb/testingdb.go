package testingdb

import (
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // register mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // register sqlite3 driver
)

const (
	mysqlTruncateTables = `
TRUNCATE tests;
`

	sqliteTruncateTables = `
DELETE FROM tests;
`
)

// MySQLDB returns a MySQL db instance for testdb.
func MySQLDB() *sqlx.DB {
	connStr := "root@tcp(localhost:3306)/tsore_development?parseTime=true"

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	}

	db, err := sqlx.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	Truncate(db)

	return db
}

// SQLiteDB returns a SQLite db instance for testdb.
func SQLiteDB(dbpath string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	Truncate(db)

	return db
}

// Truncate truncates the DB
func Truncate(db *sqlx.DB) {
	var sql []string
	switch db.DriverName() {
	case "mysql":
		sql = strings.Split(mysqlTruncateTables, "\n")
	case "sqlite3":
		sql = []string{sqliteTruncateTables}
	default:
		panic("Unknown driver")
	}

	for _, expr := range sql {
		if len(strings.TrimSpace(expr)) == 0 {
			continue
		}
		if _, err := db.Exec(expr); err != nil {
			panic(err)
		}
	}
}
