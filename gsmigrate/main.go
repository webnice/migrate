package main

import (
	"database/sql"
	"os"

	"github.com/webnice/migrate/goose"

	// Подключение драйверов баз данных.
	_ "github.com/ClickHouse/clickhouse-go" // Clickhouse
	_ "github.com/go-sql-driver/mysql"      // Mysql
	_ "github.com/lib/pq"                   // Postgres, Cockroach, Redshift
	_ "github.com/mattn/go-sqlite3"         // Sqlite
	_ "github.com/ziutek/mymysql/godrv"     // App Engine CloudSQL
)

// Точка входа в приложение.
func main() {
	var err error

	if err = execute(os.Args[1:]); err != nil {
		logFatalf("%s", err)
	}
}

// Выполнение команды утилиты.
func execute(argv []string) (err error) {
	const (
		cmdNameCreate  = "create"
		cmdUpTo        = "up-to"
		cmdDownTo      = "down-to"
		driverPostgres = "postgres"
		driverMysql    = "mysql"
	)
	var (
		cmd string
		arg *Args
		db  *sql.DB
	)

	if cmd, arg, err = args(argv); err != nil {
		return
	}
	// Команда create выполняется без подключения к базе данных.
	if cmd == cmdNameCreate {
		if err = goose.Run(cmd, nil, arg.Directory, arg.CreateName, arg.CreateType); err != nil {
			err = Errors().RunCommand.Bind(cmd, err)
		}
		return
	}
	// Проверка драйвера и установка диалекта.
	switch arg.Driver {
	case "mysql", "postgres", "sqlite3", "clickhouse":
		err = goose.SetDialect(arg.Driver)
	case "redshift":
		err = goose.SetDialect(arg.Driver)
		arg.Driver = driverPostgres
	case "tidb":
		err = goose.SetDialect(arg.Driver)
		arg.Driver = driverMysql
	default:
		err = Errors().UnknownDriver.Bind(arg.Driver)
	}
	if err != nil {
		return
	}
	// Подключение к базе данных.
	if db, err = sql.Open(arg.Driver, arg.Dsn); err != nil {
		err = Errors().ConnectDatabase.Bind(err)
		return
	}
	defer func() { _ = db.Close() }()
	db.SetConnMaxLifetime(0)
	// Выполнение команды миграции.
	switch cmd {
	case cmdUpTo, cmdDownTo:
		err = goose.Run(cmd, db, arg.Directory, arg.UpDownVersion)
	default:
		err = goose.Run(cmd, db, arg.Directory)
	}
	if err != nil {
		err = Errors().RunCommand.Bind(cmd, err)
	}

	return
}
