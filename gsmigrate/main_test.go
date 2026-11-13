package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Тестовые SQL-миграции.
const (
	testSQLBegin = `-- +goose Up
CREATE TABLE test_table_begin (id INTEGER PRIMARY KEY);
-- +goose Down
DROP TABLE test_table_begin;
`
	testSQLNext = `-- +goose Up
CREATE TABLE test_table_next (id INTEGER PRIMARY KEY);
-- +goose Down
DROP TABLE test_table_next;
`
)

// Возвращает DSN временной базы данных sqlite3.
func testTempDsn(t *testing.T) (dsn string) {
	t.Helper()
	dsn = filepath.Join(t.TempDir(), "test.db")

	return
}

// Создание временного каталога с файлами миграций.
func testMigrationDir(t *testing.T, files map[string]string) (dir string) {
	const tplWriteError = "ошибка записи файла миграции %q: %v"
	var err error

	t.Helper()
	dir = t.TempDir()
	for name, content := range files {
		if err = os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf(tplWriteError, name, err)
		}
	}

	return
}

func TestExecuteVersion(t *testing.T) {
	var err error

	if err = execute([]string{"--drv=sqlite3", "--dsn=" + testTempDsn(t), "version"}); err != nil {
		t.Errorf("не ожидалась ошибка вывода версии: %v", err)
	}
}

func TestExecuteCommands(t *testing.T) {
	const tplExecuteError = "команда %q: неожиданная ошибка: %v"
	var (
		dir string
		dsn string
		err error
	)

	dir = testMigrationDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
		"20180306100000_next.sql":  testSQLNext,
	})
	dsn = testTempDsn(t)
	// Применение всех миграций.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "up"}); err != nil {
		t.Fatalf(tplExecuteError, "up", err)
	}
	// Вывод версии и статуса.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "version"}); err != nil {
		t.Fatalf(tplExecuteError, "version", err)
	}
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "status"}); err != nil {
		t.Fatalf(tplExecuteError, "status", err)
	}
	// Применение и откат до версии.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "down-to", "20180305100000"}); err != nil {
		t.Fatalf(tplExecuteError, "down-to", err)
	}
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "up-to", "20180306100000"}); err != nil {
		t.Fatalf(tplExecuteError, "up-to", err)
	}
	// Повтор последней миграции и откат одной миграции.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "redo"}); err != nil {
		t.Fatalf(tplExecuteError, "redo", err)
	}
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "down"}); err != nil {
		t.Fatalf(tplExecuteError, "down", err)
	}
	// Откат всех миграций.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + dsn, "--dir=" + dir, "reset"}); err != nil {
		t.Fatalf(tplExecuteError, "reset", err)
	}
}

func TestExecuteCreate(t *testing.T) {
	var (
		dir     string
		matches []string
		err     error
	)

	dir = testMigrationDir(t, map[string]string{})
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + testTempDsn(t), "--dir=" + dir, "create", "add_column", "sql"}); err != nil {
		t.Fatalf("не ожидалась ошибка создания миграции: %v", err)
	}
	if matches, err = filepath.Glob(filepath.Join(dir, "*_add_column.sql")); err != nil || len(matches) != 1 {
		t.Errorf("ожидался один файл *_add_column.sql, найдено %d (ошибка: %v)", len(matches), err)
	}
	// Создание Go-миграции.
	if err = execute([]string{"--drv=sqlite3", "--dsn=" + testTempDsn(t), "--dir=" + dir, "create", "fetch_data", "go"}); err != nil {
		t.Fatalf("не ожидалась ошибка создания миграции: %v", err)
	}
	if matches, err = filepath.Glob(filepath.Join(dir, "*_fetch_data.go")); err != nil || len(matches) != 1 {
		t.Errorf("ожидался один файл *_fetch_data.go, найдено %d (ошибка: %v)", len(matches), err)
	}
}

func TestExecuteCreateError(t *testing.T) {
	var err error

	// Создание миграции в несуществующем каталоге завершается ошибкой.
	if err = execute([]string{
		"--drv=sqlite3",
		"--dsn=" + testTempDsn(t),
		"--dir=" + filepath.Join(t.TempDir(), "no_such_dir"),
		"create", "add_column", "sql",
	}); err == nil {
		t.Fatal("ожидалась ошибка создания миграции")
	}
}

func TestExecuteUnknownDriver(t *testing.T) {
	var err error

	if err = execute([]string{"--drv=oracle", "--dsn=x", "version"}); err == nil {
		t.Fatal("ожидалась ошибка неподдерживаемого драйвера")
	} else if !strings.Contains(err.Error(), "не поддерживается") {
		t.Errorf("неожиданный текст ошибки: %v", err)
	}
}

func TestExecuteBadVersion(t *testing.T) {
	var err error

	if err = execute([]string{"--drv=sqlite3", "--dsn=" + testTempDsn(t), "up-to", "abc"}); err == nil {
		t.Fatal("ожидалась ошибка нечисловой версии")
	}
}

func TestExecuteNoCommand(t *testing.T) {
	var err error

	// Отсутствие команды: kingpin завершает процесс с показом справки.
	testTerminateSwap(t)
	if err = execute([]string{"--drv=sqlite3"}); err == nil {
		t.Fatal("ожидалась ошибка отсутствия команды")
	}
}

// Тестовый случай драйвера с заведомо недоступным подключением.
var testDriverCases = []struct {
	name string
	dsn  string
}{
	{"mysql", "user:password@tcp(127.0.0.1:1)/dbname?parseTime=true"},
	{"postgres", "postgres://user:password@127.0.0.1:1/dbname?sslmode=disable"},
	{"clickhouse", "tcp://127.0.0.1:1?username=default&database=test"},
	{"redshift", "postgres://user:password@127.0.0.1:1/dbname?sslmode=disable"},
	{"tidb", "user:password@tcp(127.0.0.1:1)/dbname?parseTime=true"},
}

func TestExecuteDrivers(t *testing.T) {
	var err error

	for _, drv := range testDriverCases {
		if err = execute([]string{"--drv=" + drv.name, "--dsn=" + drv.dsn, "version"}); err == nil {
			t.Errorf("драйвер %q: ожидалась ошибка подключения", drv.name)
		}
	}
}

func TestMainFunc(t *testing.T) {
	const (
		tplFatalUnexpected = "не ожидался вывод критической ошибки: %s %v"
		tplFatalExpected   = "ожидался вывод критической ошибки"
	)
	var (
		fatalFormat string
		fatalArgs   []any
		oldArgs     []string
		oldFatal    func(string, ...any)
		oldTerm     func(int)
	)

	oldFatal = logFatalf
	logFatalf = func(format string, args ...any) { fatalFormat, fatalArgs = format, args }
	defer func() { logFatalf = oldFatal }()
	oldTerm = terminate
	terminate = func(int) {}
	defer func() { terminate = oldTerm }()
	oldArgs = os.Args
	defer func() { os.Args = oldArgs }()
	// Успешное выполнение команды.
	os.Args = []string{"gsmigrate", "--drv=sqlite3", "--dsn=" + testTempDsn(t), "version"}
	main()
	if fatalFormat != "" {
		t.Fatalf(tplFatalUnexpected, fatalFormat, fatalArgs)
	}
	// Ошибка выполнения команды.
	os.Args = []string{"gsmigrate", "--drv=unknown", "version"}
	main()
	if fatalFormat == "" {
		t.Fatal(tplFatalExpected)
	}
}
