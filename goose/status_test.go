package goose

import (
	"path/filepath"
	"testing"
)

func TestStatus(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	// Применены две из трёх миграций: одна ожидает.
	if err = UpTo(db, dir, testVersion2); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Status(db, dir); err != nil {
		t.Errorf("не ожидалась ошибка вывода статуса: %v", err)
	}
}

func TestStatusMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Status(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestStatusClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = Status(db, dir); err == nil {
		t.Fatal("ожидалась ошибка вывода статуса")
	}
}

func TestStatusQueryError(t *testing.T) {
	const (
		tplExecError = "ошибка подготовки тестовой таблицы: %v"
		errExpect    = "ожидалась ошибка чтения статуса миграции"
	)
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	// Таблица версий без колонки tstamp: запрос статуса завершается ошибкой.
	if err = createVersionTable(db); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if _, err = db.Exec("ALTER TABLE goose_db_version DROP COLUMN tstamp"); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if err = Status(db, dir); err == nil {
		t.Fatal(errExpect)
	}
}

func TestPrintMigrationStatusNoTable(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = printMigrationStatus(db, testVersion1, "20180305100000_begin.sql"); err == nil {
		t.Fatal("ожидалась ошибка чтения статуса миграции")
	}
}
