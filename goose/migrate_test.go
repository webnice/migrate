package goose

import (
	"database/sql"
	"maps"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// Тестовые версии миграций.
const (
	testVersion1 = int64(20180305100000)
	testVersion2 = int64(20180306100000)
	testVersion3 = int64(20180307100000)
	testVersion4 = int64(20180308100000)
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
	testSQLAgain = `-- +goose Up
CREATE TABLE test_table_again (id INTEGER PRIMARY KEY);
-- +goose Down
DROP TABLE test_table_again;
`
	testSQLBadDown = `-- +goose Up
CREATE TABLE test_table_baddown (id INTEGER PRIMARY KEY);
-- +goose Down
BROKEN SQL EXPRESSION;
`
	testSQLBadUp = `-- +goose Up
BROKEN SQL EXPRESSION;
-- +goose Down
DROP TABLE test_table_brokenup;
`
)

// Сохранение и восстановление глобального состояния пакета.
func testGuard(t *testing.T) {
	var (
		oldDialect    SQLDialect
		oldRegistered map[int64]*Migration
	)

	t.Helper()
	oldDialect = dialect
	oldRegistered = make(map[int64]*Migration, len(registeredGoMigrations))
	maps.Copy(oldRegistered, registeredGoMigrations)
	t.Cleanup(func() {
		dialect = oldDialect
		registeredGoMigrations = oldRegistered
	})
}

// Подключение к временной базе данных sqlite3.
func testDatabase(t *testing.T) (db *sql.DB) {
	const (
		tplDialectError = "ошибка установки диалекта sqlite3: %v"
		tplOpenError    = "ошибка подключения к тестовой базе данных: %v"
	)
	var (
		err error
		dsn string
	)

	t.Helper()
	testGuard(t)
	if err = SetDialect("sqlite3"); err != nil {
		t.Fatalf(tplDialectError, err)
	}
	dsn = filepath.Join(t.TempDir(), "test.db")
	if db, err = sql.Open("sqlite3", dsn); err != nil {
		t.Fatalf(tplOpenError, err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return
}

// Создание временного каталога с файлами миграций.
func testMigrationsDir(t *testing.T, files map[string]string) (dir string) {
	const tplWriteError = "ошибка записи файла миграции %q: %v"
	var (
		err error
	)

	t.Helper()
	dir = t.TempDir()
	for name, content := range files {
		if err = os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf(tplWriteError, name, err)
		}
	}

	return
}

// Создание каталога с тремя корректными миграциями.
func testThreeMigrationsDir(t *testing.T) (dir string) {
	return testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql":     testSQLBegin,
		"20180306100000_next.sql":      testSQLNext,
		"20180307100000_and_again.sql": testSQLAgain,
	})
}

// Диалект с заведомо ошибочным выражением записи версии миграции.
type failInsertDialect struct{ Sqlite3Dialect }

// Возвращает выражение, заведомо завершающееся ошибкой.
func (failInsertDialect) insertVersionSQL() string {
	return "INSERT INTO no_such_table (version_id) VALUES (?)"
}

// Создание тестовой миграции.
func newMigration(v int64, src string) *Migration {
	return &Migration{Version: v, Previous: -1, Next: -1, Source: src}
}

func TestMigrationSort(t *testing.T) {
	const tplSortError = "ожидался отсортированный список миграций, версия %d не на месте %d"
	var (
		ms  Migrations
		err error
	)

	ms = Migrations{
		newMigration(20129000, "20129000_and_again.sql"),
		newMigration(20120000, "20120000_begin.sql"),
		newMigration(20128000, "20128000_next.sql"),
		newMigration(20127000, "20127000_more.sql"),
	}
	if ms, err = sortAndConnectMigrations(ms); err != nil {
		t.Fatalf("не ожидалась ошибка сортировки: %v", err)
	}
	sorted := []int64{20120000, 20127000, 20128000, 20129000}
	for i, m := range ms {
		if sorted[i] != m.Version {
			t.Errorf(tplSortError, sorted[i], i)
		}
	}
	validateMigrationSort(t, ms, sorted)
}

// Проверка связей следующей и предыдущей версии.
func validateMigrationSort(t *testing.T, ms Migrations, sorted []int64) {
	t.Helper()
	for i, m := range ms {
		if sorted[i] != m.Version {
			t.Error("несортированный список миграций")
		}
		var next, prev int64

		if i == 0 {
			prev = -1
			next = ms[i+1].Version
		} else if i == len(ms)-1 {
			prev = ms[i-1].Version
			next = -1
		} else {
			prev = ms[i-1].Version
			next = ms[i+1].Version
		}
		if m.Next != next {
			t.Errorf("неверная связь Next: версия %d, получено %d, ожидалось %d", m.Version, m.Next, next)
		}
		if m.Previous != prev {
			t.Errorf("неверная связь Previous: версия %d, получено %d, ожидалось %d", m.Version, m.Previous, prev)
		}
	}
}

func TestSortAndConnectDuplicate(t *testing.T) {
	const errExpect = "ожидалась ошибка дублирующейся версии миграции"
	var (
		ms  Migrations
		err error
	)

	ms = Migrations{
		newMigration(testVersion1, "20180305100000_begin.sql"),
		newMigration(testVersion1, "20180305100000_dup.sql"),
	}
	if _, err = sortAndConnectMigrations(ms); err == nil {
		t.Fatal(errExpect)
	}
}

//goland:noinspection DuplicatedCode
func TestMigrationsNavigation(t *testing.T) {
	const (
		tplCurrentError  = "Current(%d): неожиданная ошибка: %v"
		tplCurrentExpect = "Current(%d): ожидалась ошибка ErrNoCurrentVersion"
		tplNextError     = "Next(%d): неожиданная ошибка: %v"
		tplNextExpect    = "Next(%d): ожидалась ошибка ErrNoNextVersion"
		tplPrevError     = "Previous(%d): неожиданная ошибка: %v"
		tplPrevExpect    = "Previous(%d): ожидалась ошибка ErrNoNextVersion"
		tplLastError     = "Last(): неожиданная ошибка: %v"
		tplLastExpect    = "Last(): ожидалась ошибка ErrNoNextVersion"
		tplWrongVersion  = "получена версия %d, ожидалась %d"
	)
	var (
		ms      Migrations
		err     error
		current *Migration
		next    *Migration
		prev    *Migration
		last    *Migration
	)

	ms = Migrations{
		newMigration(testVersion1, "20180305100000_begin.sql"),
		newMigration(testVersion2, "20180306100000_next.sql"),
		newMigration(testVersion3, "20180307100000_and_again.sql"),
	}
	// Поиск миграции по версии.
	if current, err = ms.Current(testVersion2); err != nil {
		t.Fatalf(tplCurrentError, testVersion2, err)
	}
	if current.Version != testVersion2 {
		t.Errorf(tplWrongVersion, current.Version, testVersion2)
	}
	if _, err = ms.Current(testVersion4); err == nil {
		t.Fatalf(tplCurrentExpect, testVersion4)
	}
	// Поиск следующей миграции.
	if next, err = ms.Next(testVersion1); err != nil {
		t.Fatalf(tplNextError, testVersion1, err)
	}
	if next.Version != testVersion2 {
		t.Errorf(tplWrongVersion, next.Version, testVersion2)
	}
	if _, err = ms.Next(testVersion3); err == nil {
		t.Fatalf(tplNextExpect, testVersion3)
	}
	// Поиск предыдущей миграции.
	if prev, err = ms.Previous(testVersion3); err != nil {
		t.Fatalf(tplPrevError, testVersion3, err)
	}
	if prev.Version != testVersion2 {
		t.Errorf(tplWrongVersion, prev.Version, testVersion2)
	}
	if _, err = ms.Previous(testVersion1); err == nil {
		t.Fatalf(tplPrevExpect, testVersion1)
	}
	// Последняя миграция.
	if last, err = ms.Last(); err != nil {
		t.Fatalf(tplLastError, err)
	}
	if last.Version != testVersion3 {
		t.Errorf(tplWrongVersion, last.Version, testVersion3)
	}
	if _, err = (Migrations{}).Last(); err == nil {
		t.Fatal(tplLastExpect)
	}
}

func TestMigrationsString(t *testing.T) {
	const (
		tplExpect = "ожидалась строка %q, получено %q"
		want      = "20180305100000_begin.sql\n20180306100000_next.sql\n"
	)
	var (
		ms Migrations
	)

	ms = Migrations{
		newMigration(testVersion1, "20180305100000_begin.sql"),
		newMigration(testVersion2, "20180306100000_next.sql"),
	}
	if ms.String() != want {
		t.Errorf(tplExpect, want, ms.String())
	}
}

func TestVersionFilter(t *testing.T) {
	const tplWrongResult = "versionFilter(%d, %d, %d) = %v, ожидалось %v"
	var tests = []struct {
		name          string
		v, cur, limit int64
		want          bool
	}{
		{"вверх в диапазоне", 5, 0, 10, true},
		{"вверх выше цели", 11, 0, 10, false},
		{"вверх ниже текущей", 1, 3, 10, false},
		{"вниз в диапазоне", 5, 10, 3, true},
		{"вниз ниже цели", 2, 10, 3, false},
		{"вниз выше текущей", 11, 10, 3, false},
		{"равенство версий", 5, 5, 5, false},
	}

	for _, tst := range tests {
		if got := versionFilter(tst.v, tst.cur, tst.limit); got != tst.want {
			t.Errorf(tplWrongResult, tst.v, tst.cur, tst.limit, got, tst.want)
		}
	}
}

func TestAddMigration(t *testing.T) {
	const (
		tplNotRegistered = "миграция %d не зарегистрирована"
		tplWrongSource   = "неправильный файл источника: %q"
	)
	var (
		mig *Migration
		ok  bool
	)

	testGuard(t)
	AddMigration(nil, nil)
	// Имя тестового файла без числового префикса — версия миграции 0.
	if mig, ok = registeredGoMigrations[0]; !ok {
		t.Fatalf(tplNotRegistered, 0)
	}
	if filepath.Base(mig.Source) != "migrate_test.go" {
		t.Errorf(tplWrongSource, mig.Source)
	}
	if !mig.Registered {
		t.Error("миграция должна быть отмечена как зарегистрированная")
	}
}

func TestAddNamedMigrationConflict(t *testing.T) {
	const errExpect = "ожидалась паника при конфликте версий миграций"

	testGuard(t)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(errExpect)
		}
	}()
	AddNamedMigration("20180305100000_begin.go", nil, nil)
	AddNamedMigration("20180305100000_other.go", nil, nil)
}

func TestCollectMigrationsNotExists(t *testing.T) {
	const errExpect = "ожидалась ошибка несуществующего каталога"
	var err error

	testGuard(t)
	if _, err = CollectMigrations(filepath.Join(t.TempDir(), "no_such_dir"), minVersion, MaxVersion); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCollectMigrations(t *testing.T) {
	const (
		tplCountError = "ожидалось %d миграций, получено %d"
		tplRegistered = "миграция %d не найдена в списке"
	)
	var (
		dir        string
		mig        *Migration
		migrations Migrations
		err        error
		ok         bool
		found      map[int64]bool
	)

	testGuard(t)
	// Go-файл миграции с зарегистрированными функциями.
	AddNamedMigration(filepath.Join("migration", "20180308100000_registered.go"), nil, nil)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql":       testSQLBegin,
		"20180306100000_next.sql":        testSQLNext,
		"20180308100000_registered.go":   "package main",
		"20180309100000_unregistered.go": "package main",
		"helpers.go":                     "package main",
	})
	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		t.Fatalf("не ожидалась ошибка сбора миграций: %v", err)
	}
	// Две SQL-миграции, зарегистрированная и незарегистрированная Go-миграции.
	if len(migrations) != 4 {
		t.Fatalf(tplCountError, 4, len(migrations))
	}
	found = make(map[int64]bool)
	for _, mig = range migrations {
		found[mig.Version] = true
	}
	for _, v := range []int64{testVersion1, testVersion2, testVersion4, 20180309100000} {
		if _, ok = found[v]; !ok {
			t.Errorf(tplRegistered, v)
		}
	}
}

func TestCollectMigrationsBadSQLName(t *testing.T) {
	const errExpect = "ожидалась ошибка некорректного имени SQL-файла"
	var err error

	testGuard(t)
	dir := testMigrationsDir(t, map[string]string{
		"withoutprefix.sql": testSQLBegin,
	})
	if _, err = CollectMigrations(dir, minVersion, MaxVersion); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCollectMigrationsDuplicateVersion(t *testing.T) {
	const errExpect = "ожидалась ошибка дублирующейся версии миграций"
	var err error

	testGuard(t)
	dir := testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
		"20180305100000_dup.sql":   testSQLNext,
	})
	if _, err = CollectMigrations(dir, minVersion, MaxVersion); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCollectMigrationsBadRegisteredName(t *testing.T) {
	const errExpect = "ожидалась ошибка некорректного имени зарегистрированной миграции"
	var err error

	testGuard(t)
	AddNamedMigration("junk.go", nil, nil)
	dir := testMigrationsDir(t, map[string]string{})
	if _, err = CollectMigrations(dir, minVersion, MaxVersion); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCollectMigrationsGlobBadPattern(t *testing.T) {
	const (
		errMkdirError        = "ошибка создания тестового каталога: %v"
		errCollectMigrations = "не ожидалась ошибка сбора миграций: %v"
		errEmpty             = "ожидался пустой список миграций, получено %d"
	)
	var (
		dir        string
		migrations Migrations
		err        error
	)

	testGuard(t)
	// Символ '[' в имени каталога делает шаблон поиска недопустимым;
	// ошибка шаблона игнорируется, миграции не найдены.
	dir = filepath.Join(t.TempDir(), "[migration")
	if err = os.Mkdir(dir, 0755); err != nil {
		t.Fatalf(errMkdirError, err)
	}
	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		t.Fatalf(errCollectMigrations, err)
	}
	if len(migrations) != 0 {
		t.Errorf(errEmpty, len(migrations))
	}
}

func TestEnsureDBVersion(t *testing.T) {
	const (
		tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
		tplInsertError  = "ошибка вставки записи версии: %v"
	)
	var (
		db      *sql.DB
		err     error
		version int64
		d       SQLDialect
	)

	db = testDatabase(t)
	d = GetDialect()
	// Чистая база данных: таблица версий создаётся автоматически.
	if version, err = EnsureDBVersion(db); err != nil || version != 0 {
		t.Errorf(tplWrongVersion, 0, version, err)
	}
	// Применённая версия и дважды откаченная версия.
	for _, row := range []struct {
		v int64
		a bool
	}{
		{testVersion2, true},
		{testVersion3, false},
		{testVersion3, false},
	} {
		if _, err = db.Exec(d.insertVersionSQL(), row.v, row.a); err != nil {
			t.Fatalf(tplInsertError, err)
		}
	}
	if version, err = EnsureDBVersion(db); err != nil || version != testVersion2 {
		t.Errorf(tplWrongVersion, testVersion2, version, err)
	}
}

func TestEnsureDBVersionNoApplied(t *testing.T) {
	const (
		tplWrongError = "ожидалась ошибка ErrNoNextVersion, получено: %v"
		tplInsertErr  = "ошибка вставки записи версии: %v"
	)
	var (
		db  *sql.DB
		err error
		d   SQLDialect
	)

	db = testDatabase(t)
	d = GetDialect()
	// Создание таблицы версий.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	// Только откаченные версии, включая начальную: текущая версия отсутствует.
	for _, row := range []struct {
		v int64
		a bool
	}{
		{0, false},
		{testVersion2, false},
	} {
		if _, err = db.Exec(d.insertVersionSQL(), row.v, row.a); err != nil {
			t.Fatalf(tplInsertErr, err)
		}
	}
	if _, err = EnsureDBVersion(db); err == nil {
		t.Fatal("ожидалась ошибка NoNextVersion")
	} else if !Errors().NoNextVersion.Is(err) {
		t.Errorf(tplWrongError, err)
	}
}

//goland:noinspection DuplicatedCode
func TestEnsureDBVersionScanError(t *testing.T) {
	const (
		tplExec      = "ошибка подготовки тестовой таблицы: %v"
		errExpect    = "ожидалась ошибка чтения записей версий"
		schemaCustom = `CREATE TABLE goose_db_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id TEXT NOT NULL,
			is_applied INTEGER NOT NULL,
			tstamp TIMESTAMP
		);`
	)
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	if _, err = db.Exec(schemaCustom); err != nil {
		t.Fatalf(tplExec, err)
	}
	if _, err = db.Exec("INSERT INTO goose_db_version (version_id, is_applied) VALUES ('abc', 1)"); err != nil {
		t.Fatalf(tplExec, err)
	}
	if _, err = EnsureDBVersion(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCreateVersionTableClosedDb(t *testing.T) {
	const errExpect = "ожидалась ошибка создания таблицы на закрытой базе данных"
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = createVersionTable(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCreateVersionTableExecError(t *testing.T) {
	const (
		tplDialectError = "ошибка установки диалекта clickhouse: %v"
		errExpect       = "ожидалась ошибка создания таблицы версий"
	)
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	if err = SetDialect("clickhouse"); err != nil {
		t.Fatalf(tplDialectError, err)
	}
	if err = createVersionTable(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCreateVersionTableInsertError(t *testing.T) {
	const errExpect = "ожидалась ошибка записи начальной версии"
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	dialect = failInsertDialect{}
	if err = createVersionTable(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestGetDBVersionError(t *testing.T) {
	const (
		tplCloseError = "ошибка закрытия базы данных: %v"
		errExpect     = "ожидалась ошибка чтения версии"
		tplWrongVer   = "ожидалась версия -1, получено %d"
	)
	var (
		db      *sql.DB
		err     error
		version int64
	)

	db = testDatabase(t)
	if err = db.Close(); err != nil {
		t.Fatalf(tplCloseError, err)
	}
	if version, err = GetDBVersion(db); err == nil {
		t.Fatal(errExpect)
	}
	if version != -1 {
		t.Errorf(tplWrongVer, version)
	}
}

func TestDbMigrationsStatus(t *testing.T) {
	const (
		tplWrongStatus = "статус версии %d: получено %v, ожидалось %v"
		tplExecError   = "ошибка вставки записи версии: %v"
	)
	var (
		db       *sql.DB
		err      error
		statuses map[int64]bool
		d        SQLDialect
	)

	db = testDatabase(t)
	d = GetDialect()
	// Отсутствующая таблица версий создаётся автоматически.
	if statuses, err = dbMigrationsStatus(db); err != nil {
		t.Fatalf("не ожидалась ошибка чтения статусов: %v", err)
	}
	if len(statuses) != 0 {
		t.Errorf(tplWrongStatus, 0, len(statuses), 0)
	}
	// Применённая и откаченная версии.
	for _, row := range []struct {
		v int64
		a bool
	}{
		{testVersion1, true},
		{testVersion2, false},
	} {
		if _, err = db.Exec(d.insertVersionSQL(), row.v, row.a); err != nil {
			t.Fatalf(tplExecError, err)
		}
	}
	if statuses, err = dbMigrationsStatus(db); err != nil {
		t.Fatalf("не ожидалась ошибка чтения статусов: %v", err)
	}
	if statuses[testVersion1] != true || statuses[testVersion2] != false {
		t.Errorf(tplWrongStatus, testVersion1, statuses[testVersion1], true)
	}
}

//goland:noinspection DuplicatedCode
func TestDbMigrationsStatusScanError(t *testing.T) {
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
	var (
		db  *sql.DB
		err error
	)

	db = testDatabase(t)
	if _, err = db.Exec(schemaCustom); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if _, err = db.Exec("INSERT INTO goose_db_version (version_id, is_applied) VALUES ('abc', 1)"); err != nil {
		t.Fatalf(tplExecError, err)
	}
	if _, err = dbMigrationsStatus(db); err == nil {
		t.Fatal(errExpect)
	}
}
