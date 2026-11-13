package goose

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestNumericComponent(t *testing.T) {
	const (
		tplWrongValue = "NumericComponent(%q) = %d, ожидалось %d"
		tplWrongError = "NumericComponent(%q), ошибка %v не соответствует справочнику %v"
		tplNoError    = "NumericComponent(%q), не ожидалась ошибка: %v"
	)
	var tests = []struct {
		name string
		in   string
		want int64
		ref  error
	}{
		{"SQL файл с путём", "/tmp/migration/20180305100000_begin.sql", testVersion1, nil},
		{"Go файл", "20180306100000_next.go", testVersion2, nil},
		{"неизвестное расширение", "20180306100000_next.txt", 0, Errors().NotMigrationFile},
		{"отсутствие расширения", "20180306100000_next", 0, Errors().NotMigrationFile},
		{"отсутствие разделителя", "20180306100000next.sql", 0, Errors().NoSeparator},
		{"нулевая версия", "0_next.sql", 0, Errors().NotPositiveVersion},
		{"отрицательная версия", "-5_next.sql", 0, Errors().NotPositiveVersion},
	}

	for _, tst := range tests[:2] {
		var (
			got int64
			err error
		)
		if got, err = NumericComponent(tst.in); err != nil {
			t.Errorf(tplNoError, tst.in, err)
			continue
		}
		if got != tst.want {
			t.Errorf(tplWrongValue, tst.in, got, tst.want)
		}
	}
	// Сравнение ошибок со справочником.
	for _, tst := range tests[2:] {
		var err error
		if _, err = NumericComponent(tst.in); !errors.Is(err, tst.ref) {
			t.Errorf(tplWrongError, tst.in, err, tst.ref)
		}
	}
	// Ошибка разбора номера версии, обёрнутая ошибка strconv.
	var err error
	if _, err = NumericComponent("abc_next.sql"); !strings.Contains(err.Error(), "Ошибка разбора номера версии") {
		t.Errorf("NumericComponent(abc_next.sql), неожиданный текст ошибки: %v", err)
	}
}

//goland:noinspection DuplicatedCode
func TestMigrationUpDownSQL(t *testing.T) {
	const (
		tplUpError      = "не ожидалась ошибка применения миграции: %v"
		tplDownError    = "не ожидалась ошибка отката миграции: %v"
		tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	)
	var (
		db      *sql.DB
		dir     string
		mig     *Migration
		err     error
		version int64
	)

	db = testDatabase(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_begin.sql": testSQLBegin,
	})
	mig = &Migration{Version: testVersion1, Next: -1, Previous: -1, Source: filepath.Join(dir, "20180305100000_begin.sql")}
	if err = mig.Up(db); err != nil {
		t.Fatalf(tplUpError, err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Fatalf(tplWrongVersion, testVersion1, version, err)
	}
	if err = mig.Down(db); err != nil {
		t.Fatalf(tplDownError, err)
	}
	if version, err = GetDBVersion(db); version != 0 {
		t.Errorf(tplWrongVersion, 0, version, err)
	}
}

func TestMigrationUpSQLFail(t *testing.T) {
	const errExpect = "ожидалась ошибка применения сломанной миграции"
	var (
		db  *sql.DB
		dir string
		mig *Migration
		err error
	)

	db = testDatabase(t)
	dir = testMigrationsDir(t, map[string]string{
		"20180305100000_broken.sql": testSQLBadDown,
	})
	mig = &Migration{Version: testVersion1, Next: -1, Previous: -1, Source: filepath.Join(dir, "20180305100000_broken.sql")}
	if err = mig.Up(db); err == nil {
		t.Fatal(errExpect)
	}
	if !strings.Contains(err.Error(), "Не удалось выполнить миграцию") {
		t.Errorf("неожиданный текст ошибки: %v", err)
	}
}

//goland:noinspection DuplicatedCode
func TestMigrationGoUpDown(t *testing.T) {
	const (
		tplUpError      = "не ожидалась ошибка применения Go-миграции: %v"
		tplDownError    = "не ожидалась ошибка отката Go-миграции: %v"
		tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
		tplNotFound     = "зарегистрированная миграция %d не найдена"
	)
	var (
		db      *sql.DB
		mig     *Migration
		err     error
		version int64
		ok      bool
	)

	db = testDatabase(t)
	testGuard(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	AddNamedMigration("20180305100000_gomig.go", func(tx *sql.Tx) (err error) {
		_, err = tx.Exec("CREATE TABLE test_go (id INTEGER PRIMARY KEY)")
		return
	}, func(tx *sql.Tx) (err error) {
		_, err = tx.Exec("DROP TABLE test_go")
		return
	})
	if mig, ok = registeredGoMigrations[testVersion1]; !ok {
		t.Fatalf(tplNotFound, testVersion1)
	}
	if err = mig.Up(db); err != nil {
		t.Fatalf(tplUpError, err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Fatalf(tplWrongVersion, testVersion1, version, err)
	}
	if err = mig.Down(db); err != nil {
		t.Fatalf(tplDownError, err)
	}
	if version, err = GetDBVersion(db); version != 0 {
		t.Errorf(tplWrongVersion, 0, version, err)
	}
}

func TestMigrationGoNotRegistered(t *testing.T) {
	const errExpect = "ожидалась ошибка незарегистрированной Go-миграции"
	var (
		db  *sql.DB
		mig *Migration
		err error
	)

	db = testDatabase(t)
	mig = &Migration{Version: testVersion1, Next: -1, Previous: -1, Source: "20180305100000_unregistered.go"}
	if err = mig.Up(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestMigrationGoFuncError(t *testing.T) {
	const errBoom = "тестовая ошибка миграции"
	var (
		db  *sql.DB
		mig *Migration
		err error
	)

	db = testDatabase(t)
	testGuard(t)
	AddNamedMigration("20180305100000_failing.go", func(tx *sql.Tx) (err error) {
		return errors.New(errBoom)
	}, nil)
	mig = registeredGoMigrations[testVersion1]
	if err = mig.Up(db); err == nil {
		t.Fatal("ожидалась ошибка выполнения функции миграции")
	}
	if !strings.Contains(err.Error(), errBoom) {
		t.Errorf("ошибка %q не содержит %q", err.Error(), errBoom)
	}
}

func TestMigrationGoNilFunc(t *testing.T) {
	const tplWrongVersion = "ожидалась версия %d, получено %d (ошибка: %v)"
	var (
		db      *sql.DB
		mig     *Migration
		err     error
		version int64
	)

	db = testDatabase(t)
	testGuard(t)
	// Создание таблицы версий перед выполнением миграций.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	AddNamedMigration("20180305100000_nilfunc.go", nil, nil)
	mig = registeredGoMigrations[testVersion1]
	if err = mig.Up(db); err != nil {
		t.Fatalf("не ожидалась ошибка миграции без функций: %v", err)
	}
	if version, err = GetDBVersion(db); version != testVersion1 {
		t.Errorf(tplWrongVersion, testVersion1, version, err)
	}
}

func TestMigrationGoBeginError(t *testing.T) {
	const errExpect = "ожидалась ошибка открытия транзакции"
	var (
		db  *sql.DB
		mig *Migration
		err error
	)

	db = testDatabase(t)
	testGuard(t)
	AddNamedMigration("20180305100000_closed.go", nil, nil)
	mig = registeredGoMigrations[testVersion1]
	if err = db.Close(); err != nil {
		t.Fatalf("ошибка закрытия базы данных: %v", err)
	}
	if err = mig.Up(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestMigrationGoInsertError(t *testing.T) {
	const errExpect = "ожидалась ошибка записи версии миграции"
	var (
		db  *sql.DB
		mig *Migration
		err error
	)

	db = testDatabase(t)
	testGuard(t)
	// Создание таблицы версий: ошибка вставки версии должна быть от диалекта.
	if _, err = EnsureDBVersion(db); err != nil {
		t.Fatalf("не ожидалась ошибка создания таблицы версий: %v", err)
	}
	AddNamedMigration("20180305100000_insertfail.go", nil, nil)
	dialect = failInsertDialect{}
	mig = registeredGoMigrations[testVersion1]
	if err = mig.Up(db); err == nil {
		t.Fatal(errExpect)
	}
}

func TestMigrationUnknownExtension(t *testing.T) {
	var (
		db  *sql.DB
		mig *Migration
		err error
	)

	db = testDatabase(t)
	mig = &Migration{Version: testVersion1, Next: -1, Previous: -1, Source: "20180305100000_notes.txt"}
	if err = mig.Up(db); err != nil {
		t.Errorf("не ожидалась ошибка миграции с неизвестным типом файла: %v", err)
	}
}
