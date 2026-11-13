package goose

import (
	"database/sql"
	"runtime"
	"strconv"
)

// Run Выполнение команды goose.
func Run(command string, db *sql.DB, dir string, argv ...string) (err error) {
	const migrationTypeGo = "go"
	var (
		version       int64
		migrationType string
	)

	switch command {
	case "up":
		err = Up(db, dir)
	case "up-by-one":
		err = UpByOne(db, dir)
	case "up-to":
		if len(argv) == 0 {
			err = Errors().UpToUsage.Bind()
			return
		}
		if version, err = strconv.ParseInt(argv[0], 10, 64); err != nil {
			err = Errors().VersionNotNumber.Bind(argv[0])
			return
		}
		err = UpTo(db, dir, version)
	case "create":
		if len(argv) == 0 {
			err = Errors().CreateUsage.Bind()
			return
		}
		migrationType = migrationTypeGo
		if len(argv) == 2 {
			migrationType = argv[1]
		}
		err = Create(db, dir, argv[0], migrationType)
	case "down":
		err = Down(db, dir)
	case "down-to":
		if len(argv) == 0 {
			err = Errors().DownToUsage.Bind()
			return
		}
		if version, err = strconv.ParseInt(argv[0], 10, 64); err != nil {
			err = Errors().VersionNotNumber.Bind(argv[0])
			return
		}
		err = DownTo(db, dir, version)
	case "redo":
		err = Redo(db, dir)
	case "reset":
		err = Reset(db, dir)
	case "status":
		err = Status(db, dir)
	case "version":
		err = Version(db, dir)
	default:
		err = Errors().UnknownCommand.Bind(command)
	}

	return
}

// AddMigration Регистрация Go-миграции, вызывающий файл определяется автоматически.
func AddMigration(up, down func(*sql.Tx) error) {
	var (
		filename string
	)

	// Получение имени файла, из которого выполняется вызов регистрации.
	_, filename, _, _ = runtime.Caller(1)
	AddNamedMigration(filename, up, down)
}

// AddNamedMigration Регистрация именованной Go-миграции.
func AddNamedMigration(filename string, up, down func(*sql.Tx) error) {
	var (
		v        int64
		mig      *Migration
		existing *Migration
		ok       bool
	)

	// Ошибка разбора невозможна, имя файла известно вызывающему коду.
	v, _ = NumericComponent(filename)
	mig = &Migration{Version: v, Next: -1, Previous: -1, Registered: true, UpFn: up, DownFn: down, Source: filename}
	if existing, ok = registeredGoMigrations[v]; ok {
		panic(Errors().MigrationConflict.Bind(filename, existing.Source))
	}
	registeredGoMigrations[v] = mig
}
