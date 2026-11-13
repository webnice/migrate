package goose

import (
	"database/sql"
	"text/template"
	"time"
)

const (
	sqlCmdPrefix        = "-- +goose "     // Префикс аннотаций goose внутри SQL-файлов миграций.
	gooseUp             = "Up"             // Аннотация начала секции прямой миграции.
	gooseDown           = "Down"           // Аннотация начала секции отката миграции.
	gooseStatementBegin = "StatementBegin" // Аннотация начала сложного SQL-выражения.
	gooseStatementEnd   = "StatementEnd"   // Аннотация окончания сложного SQL-выражения.
	gooseNoTransaction  = "NO TRANSACTION" // Аннотация отключения транзакции для файла миграции.
	extensionSQL        = ".sql"           // Расширение файла SQL-миграции.
	extensionGo         = ".go"            // Расширение файла Go-миграции.
)

var (
	MaxVersion                        = int64(9223372036854775807) // Максимально допустимая версия миграции.
	minVersion                        = int64(0)                   // Минимально допустимая версия миграции.
	registeredGoMigrations            = map[int64]*Migration{}     // Реестр Go-миграций, зарегистрированных функцией AddMigration().
	dialect                SQLDialect = &MySQLDialect{}            // Диалект базы данных, применяемый по умолчанию.
)

// SQLDialect Интерфейс абстракции специфики диалектов SQL для небольшого набора SQL-выражений goose.
type SQLDialect interface {
	createVersionTableSQL() string                // SQL-выражение создания таблицы goose_db_version.
	insertVersionSQL() string                     // SQL-выражение добавления записи в таблицу версий.
	dbVersionQuery(db *sql.DB) (*sql.Rows, error) // Запрос чтения лога версий миграций.
}

// MigrationRecord Запись лога применённых миграций.
type MigrationRecord struct {
	VersionID int64     // Номер версии миграции.
	TStamp    time.Time // Дата и время миграции.
	IsApplied bool      // Истина — миграция применена, ложь — откачена.
}

// Migration Миграция базы данных.
type Migration struct {
	Version    int64               // Номер версии миграции.
	Next       int64               // Следующая версия, или -1 если следующей версии нет.
	Previous   int64               // Предыдущая версия, или -1 если предыдущей версии нет.
	Source     string              // Путь к файлу миграции.
	Registered bool                // Истина — функции Go зарегистрированы через AddMigration().
	UpFn       func(*sql.Tx) error // Функция прямой Go-миграции.
	DownFn     func(*sql.Tx) error // Функция отката Go-миграции.
}

// Migrations Список миграций.
type Migrations []*Migration

var (
	// sqlMigrationTemplate Шаблон нового файла SQL-миграции.
	sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(
		`-- +goose Up
-- SQL в этом разделе выполняется при применении миграции.



-- +goose Down
-- SQL в этом разделе выполняется при откате миграции.

`))

	// goSQLMigrationTemplate Шаблон нового файла Go-миграции.
	goSQLMigrationTemplate = template.Must(template.New("goose.go-migration").Parse(
		`package main

import (
	"database/sql"

	"github.com/webnice/migrate/goose"
)

func init() {
	goose.AddMigration(Up{{.}}, Down{{.}})
}

// Up{{.}} Миграция применена.
func Up{{.}}(tx *sql.Tx) (err error) {
	// Этот код выполняется при применении миграции.

	return
}

// Down{{.}} Миграция откачена.
func Down{{.}}(tx *sql.Tx) (err error) {
	// Этот код выполняется при откате миграции.

	return
}
`))
)
