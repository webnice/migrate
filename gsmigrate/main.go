package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/webnice/migrate/goose"

	// Init database drivers
	_ "github.com/ClickHouse/clickhouse-go" // Clickhouse
	_ "github.com/go-sql-driver/mysql"      // Mysql
	_ "github.com/lib/pq"                   // Postgres, Cockroach, Redshift
	_ "github.com/mattn/go-sqlite3"         // Sqlite
	_ "github.com/ziutek/mymysql/godrv"     // App Engine CloudSQL
)

func main() {
	var (
		err error
		cmd string
		arg *Args
		db  *sql.DB
	)

	cmd, arg = args()
	// Exclusive execution of the create command
	if cmd == `create` {
		if err = goose.Run(cmd, nil, arg.Directory, arg.CreateName, arg.CreateType); err != nil {
			log.Fatalf("Error: %s", err)
		}
		return
	}
	// Checking driver and set dialect
	switch arg.Driver {
	case `mysql`, `postgres`, `sqlite3`, `clickhouse`:
		err = goose.SetDialect(arg.Driver)
	case `redshift`:
		err = goose.SetDialect(arg.Driver)
		arg.Driver = "postgres"
	case `tidb`:
		err = goose.SetDialect(arg.Driver)
		arg.Driver = "mysql"
	default:
		err = fmt.Errorf("%q driver not supported", arg.Driver)
	}
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	// Opening database connection
	if db, err = sql.Open(arg.Driver, arg.Dsn); err != nil {
		log.Fatalf("Connect to database error: %s", err)
	}
	db.SetConnMaxLifetime(0)
	// Running command of database migration with different arguments
	switch cmd {
	case `up-to`, `down-to`:
		err = goose.Run(cmd, db, arg.Directory, arg.UpDownVersion)
	default:
		err = goose.Run(cmd, db, arg.Directory)
	}
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
