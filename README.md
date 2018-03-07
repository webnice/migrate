# migrate

[![GoDoc](https://godoc.org/gopkg.in/webnice/migrate.v1/goose?status.svg)](https://godoc.org/gopkg.in/webnice/migrate.v1/goose)
[![Go Report Card](https://goreportcard.com/badge/github.com/webnice/migrate)](https://goreportcard.com/report/github.com/webnice/migrate)
[![Coverage Status](https://coveralls.io/repos/github/webnice/migrate/badge.svg?branch=v1)](https://coveralls.io/github/webnice/migrate?branch=v1)
[![Build Status](https://travis-ci.org/webnice/migrate.svg?branch=v1)](https://travis-ci.org/webnice/migrate)
[![CircleCI](https://circleci.com/gh/webnice/migrate/tree/v1.svg?style=svg)](https://circleci.com/gh/webnice/migrate/tree/v1)

Is a database migration tool. Manage your database schema by creating incremental SQL changes or Go functions.

Based on goose lib of `bitbucket.org/liamstask/goose`

# Supported databases

* mysql
* postgres
* cockroach
* sqlite3
* redshift
* clickhouse
* tidb


# Install

	$ go get -u gopkg.in/webnice/migrate.v1/gsmigrate

This will install the `gsmigrate` binary to your `$GOPATH/bin` directory.

# Usage

```
usage: gsmigrate [<flags>] <command> [<args> ...]

Utility for applying database migrations

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --dir="."      Directory with migration files. Default is '.'.
                 Overrides the default value for a flag from an environment variable by name 'GOOSE_DIR'
  --drv="mysql"  Driver of database. Support is [mysql, postgres, sqlite3, redshift, clickhouse, tidb].
                 Overrides the default value for a flag from an environment variable by name 'GOOSE_DRV'
  --dsn="root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true"
                 Database source name (DSN).
                 Overrides the default value for a flag from an environment variable by name 'GOOSE_DSN'

Commands:
  help [<command>...]
    Show help.

  up
    Migrate the DB to the most recent version available

  up-to [<VERSION>]
    Migrate the DB to a specific VERSION

  down
    Roll back the version by 1

  down-to [<VERSION>]
    Roll back to a specific VERSION

  redo
    Re-run the latest migration

  status
    Dump the migration status for the current DB

  version
    Print the current version of the database

  create [<NAME>] [<TYPE>]
    Creates new migration file with next version
```
## Examples:

```
    gsmigrate --drv="sqlite3" --dsn="./foo.db" status
    gsmigrate --drv="sqlite3" --dsn="./foo.db" create init sql
    gsmigrate --drv="sqlite3" --dsn="./foo.db" create add_some_column sql
    gsmigrate --drv="sqlite3" --dsn="./foo.db" create fetch_user_data go
    gsmigrate --drv="sqlite3" --dsn="./foo.db" up

    gsmigrate --drv="postgres" --dsn="user=postgres dbname=postgres sslmode=disable" status
    gsmigrate --drv="mysql" --dsn="user:password@/dbname?parseTime=true" status
    gsmigrate --drv="redshift" --dsn="postgres://user:password@qwerty.us-east-1.redshift.amazonaws.com:5439/db" status
		gsmigrate --drv="clickhouse" --dsn="tcp://localhost:9000?username=default&database=test" status
    gsmigrate --drv="tidb" --dsn="user:password@/dbname?parseTime=true" status
```
## create

Create a new SQL migration

    $ gsmigrate create add_some_column sql
    $ Created new file: YYYYMMDDhhmmss_add_some_column.sql

Edit the newly created file to define the behavior of your migration.

You can also create a Go migration, if you then invoke it with your own goose binary:

    $ gsmigrate create fetch_user_data go
    $ Created new file: YYYYMMDDhhmmss_fetch_user_data.go

## up

Apply all available migrations

    $ gsmigrate up
    $ OK    20180305100000_begin.sql
    $ OK    20180306100000_next.sql
    $ OK    20180307100000_and_again.sql
		$ goose: no migrations to run. current version: 20180307100000

## up-to

Migrate up to a specific version

    $ gsmigrate up-to 20180306100000
    $ OK    20180305100000_begin.sql
		$ OK    20180306100000_next.sql
		$ goose: no migrations to run. current version: 20180306100000

## down

Roll back a single migration from the current version

    $ gsmigrate down
    $ OK    20180307100000_and_again.sql
		$ goose: no migrations to run. current version: 20180306100000

## down-to

Roll back migrations to a specific version

    $ gsmigrate down-to 20180305100000
    $ OK    20180307100000_and_again.sql
		$ OK    20180306100000_next.sql
		$ goose: no migrations to run. current version: 20180305100000

## redo

Roll back the most recently applied migration, then run it again

    $ gsmigrate redo
		$ -- SQL in this section is executed when the migration is rolled back..
    $ OK    20180307100000_and_again.sql
		$ -- SQL in this section is executed when the migration is applied..
		$ OK    20180307100000_and_again.sql

## status

Print the status of all migrations:

    $ gsmigrate status
    $ Applied At                  Migration
		$ =======================================
		$ Wed Mar  7 14:57:35 2018 -- 20180305100000_begin.sql
		$ Wed Mar  7 15:01:20 2018 -- 20180306100000_next.sql
		$ Pending                  -- 20180307100000_and_again.sql

## version

Print the current version of the database:

    $ gsmigrate version
    $ goose: version 20180306100000

# Migrations

gsmigrate supports migrations written in SQL or in Go.

## SQL Migrations

A sample SQL migration looks like:

```sql
-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

Notice the annotations in the comments. Any statements following `-- +goose Up` will be executed as part of a forward migration, and any statements following `-- +goose Down` will be executed as part of a rollback.

By default, all migrations are run within a transaction. Some statements like `CREATE DATABASE`, however, cannot be run within a transaction. You may optionally add `-- +goose NO TRANSACTION` to the top of your migration
file in order to skip transactions within that specific migration file. Both Up and Down migrations within this file will be run without transactions.

By default, SQL statements are delimited by semicolons - in fact, query statements must end with a semicolon to be properly recognized by goose.

More complex statements (PL/pgSQL) that have semicolons within them must be annotated with `-- +goose StatementBegin` and `-- +goose StatementEnd` to be properly recognized. For example:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation( DATE, DATE )
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR( d, 'YYYY_MM' )
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR( d, 'YYYY-MM-DD 00:00:00' )
      || ''' AND created_at < timestamp '''
      || TO_CHAR( d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00' )
      || ''' ) ) inherits ( histories );'
    FROM generate_series( $1, $2, '1 month' ) AS d
  LOOP
    EXECUTE create_query;
  END LOOP;  -- LOOP END
END;         -- FUNCTION END
$$
language plpgsql;
-- +goose StatementEnd
```

## Go Migrations

1. Create your own goose binary
2. Import `gopkg.in/webnice/migrate.v1/goose`
3. Register your migration functions
4. Run goose command, ie. `goose.Up(db *sql.DB, dir string)`

A sample Go migration looks like:

```go
package main

import (
	"database/sql"

	"gopkg.in/webnice/migrate.v1/goose"
)

func init() {
	goose.AddMigration(Up20180307100000, Down20180307100000)
}

// Up20180307100000 Migration applied
func Up20180307100000(tx *sql.Tx) (err error) {
	// This code is executed when the migration is applied.
	return
}

// Down20180307100000 Migration rolled back
func Down20180307100000(tx *sql.Tx) (err error) {
	// This code is executed when the migration is rolled back.
	return
}
```

## License

Licensed under [MIT License](./LICENSE)
