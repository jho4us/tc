# testdb usage

## Setup/Migration

This directory stores [goose](https://bitbucket.org/liamstask/goose/) db migration scripts for various DB backends.
Currently supported:
 - MySQL in mysql
 - SQLite in sqlite

### Get goose

    go get bitbucket.org/liamstask/goose/cmd/goose

### Use goose to start and terminate a MySQL DB
To start a MySQL using goose:

    goose -path $GOPATH/src/ispringsolutions/tc/testdb/mysql up

To tear down a MySQL DB using goose

    goose -path $GOPATH/src/ispringsolutions/tc/testdb/mysql down

Note: the administration of MySQL DB is not included. We assume
the databases being connected to are already created and access control
is properly handled.

### Use goose to start and terminate a SQLite DB
To start a SQLite DB using goose:

    goose -path $GOPATH/src/ispringsolutions/tc/testdb/sqlite up

To tear down a SQLite DB using goose

    goose -path $GOPATH/src/ispringsolutions/tc/testdb/sqlite down

## DB Configuration

We're take a -db-config flag. Create a file with a
JSON dictionary:

    {"driver":"sqlite3","data_source":"tests.db"}

or

    {"driver":"mysql","data_source":"user:password@tcp(hostname:3306)/db?parseTime=true"}
