package goose

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// Список всех диалектов баз данных.
var testDialects = []SQLDialect{
	&PostgresDialect{},
	&MySQLDialect{},
	&Sqlite3Dialect{},
	&RedshiftDialect{},
	&TiDBDialect{},
	&ClickHouseDialect{},
}

func TestSetDialect(t *testing.T) {
	const (
		tplTypeError  = "SetDialect(%q): ожидался тип %T, получено %T"
		errExpect     = "ожидалась ошибка неизвестного диалекта"
		tplWrongError = "ошибка %q не содержит %q"
	)
	var names = []string{"postgres", "mysql", "sqlite3", "redshift", "tidb", "clickhouse"}

	for i, name := range names {
		if err := SetDialect(name); err != nil {
			t.Fatalf("не ожидалась ошибка установки диалекта %q: %v", name, err)
		}
		if fmt.Sprintf("%T", GetDialect()) != fmt.Sprintf("%T", testDialects[i]) {
			t.Errorf(tplTypeError, name, testDialects[i], GetDialect())
		}
	}
	if err := SetDialect("oracle"); err == nil {
		t.Fatal(errExpect)
	} else if !strings.Contains(err.Error(), "неизвестный диалект") {
		t.Errorf(tplWrongError, err.Error(), "неизвестный диалект")
	}
}

func TestDialectSQLStrings(t *testing.T) {
	const (
		tplNoTable  = "диалект %T: выражение создания таблицы не содержит goose_db_version"
		tplNoInsert = "диалект %T: пустое выражение записи версии"
	)

	for _, d := range testDialects {
		if !strings.Contains(d.createVersionTableSQL(), "goose_db_version") {
			t.Errorf(tplNoTable, d)
		}
		if d.insertVersionSQL() == "" {
			t.Errorf(tplNoInsert, d)
		}
	}
}

func TestDialectVersionQuery(t *testing.T) {
	const (
		tplDialectError = "ошибка установки диалекта %q: %v"
		tplQueryError   = "диалект %T: неожиданная ошибка запроса версий: %v"
		errExpect       = "диалект %T: ожидалась ошибка запроса версий"
	)
	var (
		db     *sql.DB
		rows   *sql.Rows
		err    error
		schema string
	)

	db = testDatabase(t)
	schema = Sqlite3Dialect{}.createVersionTableSQL()
	if _, err = db.Exec(schema); err != nil {
		t.Fatalf("ошибка создания таблицы версий: %v", err)
	}
	for _, name := range []string{"postgres", "mysql", "sqlite3", "redshift", "tidb", "clickhouse"} {
		if err = SetDialect(name); err != nil {
			t.Fatalf(tplDialectError, name, err)
		}
		if rows, err = GetDialect().dbVersionQuery(db); err != nil {
			t.Errorf(tplQueryError, GetDialect(), err)
			continue
		}
		_ = rows.Close()
	}
	// Запрос к базе данных без таблицы версий.
	db = testDatabase(t)
	for _, name := range []string{"postgres", "mysql", "sqlite3", "redshift", "tidb", "clickhouse"} {
		if err = SetDialect(name); err != nil {
			t.Fatalf(tplDialectError, name, err)
		}
		if _, err = GetDialect().dbVersionQuery(db); err == nil {
			t.Errorf(errExpect, GetDialect())
		}
	}
}
