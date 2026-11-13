package goose

import (
	"path/filepath"
	"testing"
)

func TestRedo(t *testing.T) {
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
	if err = Redo(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка повторного выполнения миграции: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion3 {
		t.Errorf(tplWrongVersion, testVersion3, version, err)
	}
}

func TestRedoNoMigration(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = Redo(db, dir); err == nil {
		t.Fatal("ожидалась ошибка повторного выполнения без применённых миграций")
	}
}

func TestRedoClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = Redo(db, dir); err == nil {
		t.Fatal("ожидалась ошибка чтения версии базы данных")
	}
}

func TestRedoMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Redo(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestRedoBadDown(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_baddown.sql": testSQLBadDown,
	})
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Redo(db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката сломанной миграции")
	}
}
