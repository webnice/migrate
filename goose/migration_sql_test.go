package goose

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// Тестовые скрипты миграций.
const (
	testSQLNoTx = `-- +goose NO TRANSACTION
-- +goose Up
CREATE TABLE test_table_notx (id INTEGER PRIMARY KEY);
-- +goose Down
DROP TABLE test_table_notx;
`
	testSQLUnfinished = `-- +goose Up
CREATE TABLE test_unfinished (id INTEGER PRIMARY KEY)
`
	testSQLNoAnnotations = `CREATE TABLE test_noannotations (id INTEGER PRIMARY KEY);
`
	testSQLUnclosedStatement = `-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_unclosed (id INTEGER PRIMARY KEY);
`
	testSQLBroken = `-- +goose Up
CREATE TABLE test_broken (id INTEGER PRIMARY KEY);
BROKEN SQL EXPRESSION;
DROP TABLE test_broken;
-- +goose Down
DROP TABLE test_broken;
`
	testSQLNoTxBroken = `-- +goose NO TRANSACTION
-- +goose Up
CREATE TABLE test_notx_broken (id INTEGER PRIMARY KEY);
BROKEN SQL EXPRESSION;
DROP TABLE test_notx_broken;
-- +goose Down
DROP TABLE test_notx_broken;
`
)

func TestSemicolons(t *testing.T) {
	const tplWrongResult = "endsWithSemicolon(%q) = %v, ожидалось %v"
	var tests = []struct {
		line   string
		result bool
	}{
		{"END;", true},
		{"END; -- comment", true},
		{"END   ; -- comment", true},
		{"END -- comment", false},
		{"END -- comment ;", false},
		{"END \" ; \" -- comment", false},
		{"", false},
	}

	for _, tst := range tests {
		if got := endsWithSemicolon(tst.line); got != tst.result {
			t.Errorf(tplWrongResult, tst.line, got, tst.result)
		}
	}
}

//goland:noinspection DuplicatedCode
func TestGetSQLStatements(t *testing.T) {
	const (
		tplWrongCount = "getSQLStatements, получено %d выражений, ожидалось %d"
		errExpect     = "ожидалась ошибка отсутствия аннотаций"
	)
	var (
		stmts []string
		tx    bool
		err   error
	)

	// Прямая миграция сложного скрипта.
	if stmts, tx, err = getSQLStatements(strings.NewReader(funCTxt), true); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 2 || !tx {
		t.Errorf(tplWrongCount, len(stmts), 2)
	}
	// Откат сложного скрипта.
	if stmts, _, err = getSQLStatements(strings.NewReader(funCTxt), false); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 2 {
		t.Errorf(tplWrongCount, len(stmts), 2)
	}
	// Скрипт с несколькими секциями Up/Down.
	if stmts, _, err = getSQLStatements(strings.NewReader(multiTxt), true); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 2 {
		t.Errorf(tplWrongCount, len(stmts), 2)
	}
	// Откат скрипта с несколькими секциями Up/Down.
	if stmts, _, err = getSQLStatements(strings.NewReader(multiTxt), false); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 2 {
		t.Errorf(tplWrongCount, len(stmts), 2)
	}
	// Отключение транзакции.
	if _, tx, err = getSQLStatements(strings.NewReader(testSQLNoTx), true); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if tx {
		t.Error("ожидалось отключение транзакции")
	}
	// Незавершённое выражение: предупреждение без ошибки.
	if stmts, _, err = getSQLStatements(strings.NewReader(testSQLUnfinished), true); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 0 {
		t.Errorf(tplWrongCount, len(stmts), 0)
	}
	// Незакрытый блок StatementBegin: предупреждение без ошибки.
	if stmts, _, err = getSQLStatements(strings.NewReader(testSQLUnclosedStatement), true); err != nil {
		t.Fatalf("не ожидалась ошибка разбора: %v", err)
	}
	if len(stmts) != 0 {
		t.Errorf(tplWrongCount, len(stmts), 0)
	}
	// Отсутствие аннотаций: ошибка.
	if _, _, err = getSQLStatements(strings.NewReader(testSQLNoAnnotations), true); err == nil {
		t.Fatal(errExpect)
	}
}

// Читатель, завершающийся ошибкой.
type testErrorReader struct{}

// Read Всегда возвращает ошибку.
func (testErrorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("тестовая ошибка чтения")
}

func TestGetSQLStatementsReadError(t *testing.T) {
	var err error

	if _, _, err = getSQLStatements(testErrorReader{}, true); err == nil {
		t.Fatal("ожидалась ошибка чтения файла миграции")
	}
}

func TestRunSQLMigrationSuccess(t *testing.T) {
	const (
		tplRunError     = "не ожидалась ошибка выполнения миграции: %v"
		tplWrongCount   = "ожидалось %d записей в таблице, получено %d"
		tplQueryError   = "ошибка проверки таблицы: %v"
		tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	)
	var (
		db      *sql.DB
		dir     string
		err     error
		version int64
		count   int
	)

	db = testDatabase(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
	})
	// Применение миграции.
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_begin.sql"), testVersion1, true); err != nil {
		t.Fatalf(tplRunError, err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Fatalf(tplWrongVersion, testVersion1, version, err)
	}
	// Откат миграции.
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_begin.sql"), testVersion1, false); err != nil {
		t.Fatalf(tplRunError, err)
	}
	if err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE name='test_table_begin'").Scan(&count); err != nil {
		t.Fatalf(tplQueryError, err)
	}
	if count != 0 {
		t.Errorf(tplWrongCount, 0, count)
	}
}

func TestRunSQLMigrationNoTx(t *testing.T) {
	const tplRunError = "не ожидалась ошибка выполнения миграции без транзакции: %v"
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_notx.sql": testSQLNoTx,
	})
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_notx.sql"), testVersion1, true); err != nil {
		t.Fatalf(tplRunError, err)
	}
}

func TestRunSQLMigrationOpenError(t *testing.T) {
	const errExpect = "ожидалась ошибка открытия файла миграции"
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	if err = runSQLMigration(db, filepath.Join(t.TempDir(), "no_such_file.sql"), testVersion1, true); err == nil {
		t.Fatal(errExpect)
	}
}

func TestRunSQLMigrationBadStatement(t *testing.T) {
	const (
		tplRunError   = "не ожидалась ошибка проверки таблицы: %v"
		tplWrongCount = "ожидалось %d записей в таблице, получено %d"
	)
	var (
		db    *sql.DB
		dir   string
		err   error
		count int
	)

	db = testDatabase(t)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_broken.sql": testSQLBroken,
	})
	// Ошибка выражения: миграция откатывается транзакцией.
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_broken.sql"), testVersion1, true); err == nil {
		t.Fatal("ожидалась ошибка выполнения выражения")
	}
	if err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE name='test_broken'").Scan(&count); err != nil {
		t.Fatalf(tplRunError, err)
	}
	if count != 0 {
		t.Errorf(tplWrongCount, 0, count)
	}
}

func TestRunSQLMigrationNoTxBadStatement(t *testing.T) {
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_notx_broken.sql": testSQLNoTxBroken,
	})
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_notx_broken.sql"), testVersion1, true); err == nil {
		t.Fatal("ожидалась ошибка выполнения выражения без транзакции")
	}
}

func TestRunSQLMigrationNoAnnotations(t *testing.T) {
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_noanno.sql": testSQLNoAnnotations,
	})
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_noanno.sql"), testVersion1, true); err == nil {
		t.Fatal("ожидалась ошибка отсутствия аннотаций")
	}
}

func TestRunSQLMigrationBeginError(t *testing.T) {
	const errExpect = "ожидалась ошибка открытия транзакции"
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
	})
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_begin.sql"), testVersion1, true); err == nil {
		t.Fatal(errExpect)
	}
}

func TestRunSQLMigrationInsertErrorTx(t *testing.T) {
	const errExpect = "ожидалась ошибка записи версии миграции"
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	testGuard(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	dialect = failInsertDialect{}
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
	})
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_begin.sql"), testVersion1, true); err == nil {
		t.Fatal(errExpect)
	}
}

func TestRunSQLMigrationInsertErrorNoTx(t *testing.T) {
	const errExpect = "ожидалась ошибка записи версии миграции"
	var (
		db  *sql.DB
		dir string
		err error
	)

	db = testDatabase(t)
	testGuard(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	dialect = failInsertDialect{}
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_notx.sql": testSQLNoTx,
	})
	if err = runSQLMigration(db, filepath.Join(dir, "20180305100000_notx.sql"), testVersion1, true); err == nil {
		t.Fatal(errExpect)
	}
}

var funCTxt = `-- +goose Up
CREATE TABLE IF NOT EXISTS histories (
  id                BIGSERIAL  PRIMARY KEY,
  current_value     varchar(2000) NOT NULL,
  created_at      timestamp with time zone  NOT NULL
);

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

-- +goose Down
drop function histories_partition_creation(DATE, DATE);
drop TABLE histories;
`

// Скрипт с несколькими переходами Up/Down.
var multiTxt = `-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;

-- +goose Up
CREATE TABLE fancier_post (
    id int NOT NULL,
    title text,
    body text,
    created_on timestamp without time zone,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE fancier_post;
`
