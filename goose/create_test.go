package goose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestCreateSQL(t *testing.T) {
	const (
		tplCreateError  = "не ожидалась ошибка создания миграции: %v"
		tplReadError    = "ошибка чтения созданного файла: %v"
		tplWrongContent = "содержимое файла не содержит аннотацию %q"
	)
	var (
		dir     string
		content []byte
		err     error
	)

	dir = testMigrationsDir(t, map[string]string{})
	if err = Create(nil, dir, "init", "sql"); err != nil {
		t.Fatalf(tplCreateError, err)
	}
	// Поиск созданного файла по шаблону имени.
	matches, _ := filepath.Glob(filepath.Join(dir, "*_init.sql"))
	if len(matches) != 1 {
		t.Fatalf("ожидался один созданный файл *_init.sql, найдено %d", len(matches))
	}
	if content, err = os.ReadFile(matches[0]); err != nil {
		t.Fatalf(tplReadError, err)
	}
	if !strings.Contains(string(content), "-- +goose Up") {
		t.Errorf(tplWrongContent, "-- +goose Up")
	}
}

func TestCreateGo(t *testing.T) {
	const (
		errCreateError  = "не ожидалась ошибка создания миграции: %v"
		errGlob         = "ожидался один созданный файл *_fetch_data.go, найдено %d"
		errReadError    = "ошибка чтения созданного файла: %v"
		errWrongContent = "содержимое Go-файла не содержит %q"
	)
	var (
		dir     string
		content []byte
		err     error
	)

	dir = testMigrationsDir(t, map[string]string{})
	if err = Create(nil, dir, "fetch_data", "go"); err != nil {
		t.Fatalf(errCreateError, err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "*_fetch_data.go"))
	if len(matches) != 1 {
		t.Fatalf(errGlob, len(matches))
	}
	if content, err = os.ReadFile(matches[0]); err != nil {
		t.Fatalf(errReadError, err)
	}
	if !strings.Contains(string(content), "package main") {
		t.Errorf(errWrongContent, "package main")
	}
}

func TestCreateWithTemplate(t *testing.T) {
	const (
		errCustom             = "-- тестовый шаблон версии {{.}}\n"
		errCreateWithTemplate = "не ожидалась ошибка создания миграции по шаблону: %v"
		errGlob               = "ожидался один созданный файл *_custom.sql, найдено %d"
		errReadErr            = "ошибка чтения созданного файла: %v"
		errExpect             = "содержимое файла не содержит строку шаблона"
	)
	var (
		dir     string
		tmpl    *template.Template
		content []byte
		err     error
	)

	dir = testMigrationsDir(t, map[string]string{})
	tmpl = template.Must(template.New("custom").Parse(errCustom))
	if err = CreateWithTemplate(nil, dir, tmpl, "custom", "sql"); err != nil {
		t.Fatalf(errCreateWithTemplate, err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "*_custom.sql"))
	if len(matches) != 1 {
		t.Fatalf(errGlob, len(matches))
	}
	if content, err = os.ReadFile(matches[0]); err != nil {
		t.Fatalf(errReadErr, err)
	}
	if !strings.Contains(string(content), "тестовый шаблон версии") {
		t.Error(errExpect)
	}
}

func TestCreateWithTemplateExecuteError(t *testing.T) {
	const errExpect = "ожидалась ошибка выполнения шаблона"
	var (
		dir  string
		tmpl *template.Template
		err  error
	)

	dir = testMigrationsDir(t, map[string]string{})
	tmpl = template.Must(template.New("bad").Parse("{{.NoSuchField}}"))
	if err = CreateWithTemplate(nil, dir, tmpl, "bad", "sql"); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCreateDirectoryNotExists(t *testing.T) {
	const errExpect = "ожидалась ошибка создания миграции в несуществующем каталоге"
	var err error

	if err = Create(nil, filepath.Join(t.TempDir(), "no_such_dir"), "init", "sql"); err == nil {
		t.Fatal(errExpect)
	}
}

func TestCreateVersionCollision(t *testing.T) {
	const (
		fmtVersionLayout    = "20060102150405"
		errWriteError       = "ошибка записи файла миграции: %v"
		errCreateErr        = "не ожидалась ошибка создания миграции: %v"
		errCollision        = "не удалось воспроизвести коллизию версий за %d попыток"
		errNumericComponent = "ошибка разбора версии файла %q: %v"
		errVersions         = "ожидались версии %d и %d, получено %d и %d"
	)
	var (
		dir     string
		files   []string
		version string
		first   int64
		second  int64
		err     error
	)

	// Коллизия версий зависит от времени — выполняется несколько попыток.
	for range 5 {
		dir = testMigrationsDir(t, map[string]string{})
		version = time.Now().In(time.UTC).Format(fmtVersionLayout)
		if err = os.WriteFile(filepath.Join(dir, version+"_collision.sql"), []byte(testSQLBegin), 0644); err != nil {
			t.Fatalf(errWriteError, err)
		}
		if err = Create(nil, dir, "collision", "sql"); err != nil {
			t.Fatalf(errCreateErr, err)
		}
		if files, err = filepath.Glob(filepath.Join(dir, "*_collision.sql")); err != nil {
			t.Fatalf(errCreateErr, err)
		}
		if len(files) == 2 {
			break
		}
	}
	if len(files) != 2 {
		t.Fatalf(errCollision, 5)
	}
	if first, err = NumericComponent(files[0]); err != nil {
		t.Fatalf(errNumericComponent, files[0], err)
	}
	if second, err = NumericComponent(files[1]); err != nil {
		t.Fatalf(errNumericComponent, files[1], err)
	}
	// Вторая версия должна быть на единицу больше первой.
	if second != first+1 && first != second+1 {
		t.Errorf(errVersions, first, first+1, first, second)
	}
}

func TestWriteTemplateToFileExists(t *testing.T) {
	const errExpect = "ожидалась ошибка создания существующего файла"
	var (
		dir  string
		path string
		err  error
	)

	dir = testMigrationsDir(t, map[string]string{"20180305100000_begin.sql": testSQLBegin})
	path = filepath.Join(dir, "20180305100000_begin.sql")
	if _, err = writeTemplateToFile(path, sqlMigrationTemplate, "20180305100000"); err == nil {
		t.Fatal(errExpect)
	}
}

func TestWriteTemplateToFileCreateError(t *testing.T) {
	const errExpect = "ожидалась ошибка создания файла в несуществующем каталоге"
	var err error

	path := filepath.Join(t.TempDir(), "no_such_dir", "20180305100000_x.sql")
	if _, err = writeTemplateToFile(path, sqlMigrationTemplate, "20180305100000"); err == nil {
		t.Fatal(errExpect)
	}
}

func TestWriteTemplateToFileExecuteError(t *testing.T) {
	const errExpect = "ожидалась ошибка выполнения шаблона"
	var (
		dir  string
		tmpl *template.Template
		err  error
	)

	dir = testMigrationsDir(t, map[string]string{})
	tmpl = template.Must(template.New("bad").Parse("{{.NoSuchField}}"))
	if _, err = writeTemplateToFile(filepath.Join(dir, "20180305100000_bad.sql"), tmpl, "20180305100000"); err == nil {
		t.Fatal(errExpect)
	}
}
