package goose

import (
	"path/filepath"
	"testing"
)

func TestUp(t *testing.T) {
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
	if version, err = GetDBVersion(db); version != testVersion3 {
		t.Errorf(tplWrongVersion, testVersion3, version, err)
	}
}

func TestUpTo(t *testing.T) {
	const tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	var (
		db      = testDatabase(t)
		dir     = testThreeMigrationsDir(t)
		version int64
		err     error
	)

	if err = UpTo(db, dir, testVersion2); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграций: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion2 {
		t.Fatalf(tplWrongVersion, testVersion2, version, err)
	}
	// Целевая версия ниже текущей: миграций нет.
	if err = UpTo(db, dir, testVersion1); err != nil {
		t.Errorf("не ожидалась ошибка применения миграций: %v", err)
	}
}

func TestUpByOne(t *testing.T) {
	const tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	var (
		db      = testDatabase(t)
		dir     = testThreeMigrationsDir(t)
		version int64
		err     error
	)

	if err = UpByOne(db, dir); err != nil {
		t.Fatalf("не ожидалась ошибка применения миграции: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Fatalf(tplWrongVersion, testVersion1, version, err)
	}
	// Применять больше нечего: ошибки нет.
	if err = UpByOne(db, dir); err != nil {
		t.Errorf("не ожидалась ошибка применения миграции: %v", err)
	}
}

func TestUpByOneMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = UpByOne(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestUpByOneClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = UpByOne(db, dir); err == nil {
		t.Fatal("ожидалась ошибка чтения версии базы данных")
	}
}

func TestUpMissingDir(t *testing.T) {
	var err error

	db := testDatabase(t)
	if err = Up(db, filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Fatal("ожидалась ошибка несуществующего каталога")
	}
}

func TestUpBadMigration(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql":  testSQLBegin,
		"20180306100000_broken.sql": testSQLBadUp,
	})
	if err = Up(db, dir); err == nil {
		t.Fatal("ожидалась ошибка применения сломанной миграции")
	}
}

func TestUpClosedDb(t *testing.T) {
	var err error

	db := testDatabase(t)
	dir := testThreeMigrationsDir(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = Up(db, dir); err == nil {
		t.Fatal("ожидалась ошибка чтения версии базы данных")
	}
}
