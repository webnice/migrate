package goose

import (
	"path/filepath"
	"testing"
)

func TestDown(t *testing.T) {
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
	if err = Down(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка отката миграции: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion2 {
		t.Errorf(tplWrongVersion, testVersion2, version, err)
	}
}

func TestDownNoMigration(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = Down(db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката без применённых миграций")
	}
}

func TestDownTo(t *testing.T) {
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
	if err = DownTo(db, dir, testVersion1); err != nil {
		t.Fatalf("не ожидалась ошибка отката миграций: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Errorf(tplWrongVersion, testVersion1, version, err)
	}
	// Откат до нулевой версии.
	if err = DownTo(db, dir, 0); err != nil {
		t.Errorf("не ожидалась ошибка отката миграций: %v", err)
	}
	if version, err = GetDBVersion(db); version != 0 {
		t.Errorf(tplWrongVersion, 0, version, err)
	}
}

func TestDownMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Down(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestDownClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = Down(db, dir); err == nil {
		t.Fatal("ожидалась ошибка чтения версии базы данных")
	}
}

func TestDownToMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = DownTo(db, filepath.Join(t.TempDir(), "no_such_dir"), 0); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestDownToBadMigration(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql":   testSQLBegin,
		"20180306100000_baddown.sql": testSQLBadDown,
	})
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = DownTo(db, dir, 0); err == nil {
		t.Fatal("ожидалась ошибка отката сломанной миграции")
	}
}

func TestDownToClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = DownTo(db, dir, testVersion1); err == nil {
		t.Fatal("ожидалась ошибка чтения версии базы данных")
	}
}

func TestDownBadMigration(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_baddown.sql": testSQLBadDown,
	})
	if err = Up(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if err = Down(db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката сломанной миграции")
	}
}
