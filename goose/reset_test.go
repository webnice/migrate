package goose

import (
	"path/filepath"
	"testing"
)

func TestReset(t *testing.T) {
	const tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	var (
		db      = testDatabase(t)
		dir     = testThreeMigrationsDir(t)
		version int64
		err     error
	)

	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Reset(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка отката миграций: %v", err)
	}
	if version, err = GetDBVersion(db); version != 0 {
		t.Errorf(tplWrongVersion, 0, version, err)
	}
}

func TestResetPartial(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	// Применена только первая миграция: остальные пропускаются.
	if err = UpTo(db, dir, testVersion1); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Reset(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка отката миграций: %v", err)
	}
}

func TestResetMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Reset(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestResetScanError(t *testing.T) {
	const (
		tplExecError = "ошибка подготовки тестовой таблицы: %v"
		errExpect    = "ожидалась ошибка чтения статусов миграций"
		schemaCustom = `CREATE TABLE goose_db_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id TEXT NOT NULL,
			is_applied INTEGER NOT NULL,
			tstamp TIMESTAMP
		);`
	)
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if _, err = db.Exec(schemaCustom); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if _, err = db.Exec("INSERT INTO goose_db_version (version_id, is_applied) VALUES ('abc', 1)"); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if err = Reset(db, dir); err == nil {
		t.Fatal(errExpect)
	}
}

func TestResetBreakOnError(t *testing.T) {
	var err error

	db := testDatabase(t)
	// Вторая (старшая) миграция падает при откате: после ошибки цикл прерывается.
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql":     testSQLBegin,
		"20180306100000_baddown.sql":   testSQLBadDown,
		"20180307100000_and_again.sql": testSQLAgain,
	})
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Reset(db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката сломанной миграции")
	}
}

func TestResetBadDown(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_baddown.sql": testSQLBadDown,
	})
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Reset(db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката сломанной миграции")
	}
}
