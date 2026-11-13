package goose

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommands(t *testing.T) {
	const tplRunError = "команда %q, не ожидалась ошибка: %v"
	var (
		db  = testDatabase(t)
		dir = testThreeMigrationsDir(t)
		err error
	)

	// Применение всех миграций.
	if err = Run("up", db, dir); err != nil {
		t.Fatalf(tplRunError, "up", err)
	}
	// Вывод версии и статуса.
	if err = Run("version", db, dir); err != nil {
		t.Fatalf(tplRunError, "version", err)
	}
	if err = Run("status", db, dir); err != nil {
		t.Fatalf(tplRunError, "status", err)
	}
	// Применение до уже достигнутой версии, миграций нет.
	if err = Run("up-to", db, dir, "20180306100000"); err != nil {
		t.Fatalf(tplRunError, "up-to", err)
	}
	// Одна миграция, применять нечего.
	if err = Run("up-by-one", db, dir); err != nil {
		t.Fatalf(tplRunError, "up-by-one", err)
	}
	// Повтор последней миграции.
	if err = Run("redo", db, dir); err != nil {
		t.Fatalf(tplRunError, "redo", err)
	}
	// Откат до версии.
	if err = Run("down-to", db, dir, "20180305100000"); err != nil {
		t.Fatalf(tplRunError, "down-to", err)
	}
	// Откат одной миграции.
	if err = Run("down", db, dir); err != nil {
		t.Fatalf(tplRunError, "down", err)
	}
	// Откат при отсутствии применённых миграций - ошибка.
	if err = Run("down", db, dir); err == nil {
		t.Fatal("ожидалась ошибка отката без применённых миграций")
	} else if !strings.Contains(err.Error(), "Миграция 0 не найдена") {
		t.Errorf("неожиданный текст ошибки: %v", err)
	}
	// Откат всех миграций.
	if err = Run("up", db, dir); err != nil {
		t.Fatalf(tplRunError, "up", err)
	}
	if err = Run("reset", db, dir); err != nil {
		t.Fatalf(tplRunError, "reset", err)
	}
	// Создание файлов миграций без подключения к базе данных.
	if err = Run("create", nil, dir, "add_column", "sql"); err != nil {
		t.Fatalf(tplRunError, "create", err)
	}
	if err = Run("create", nil, dir, "fetch_data"); err != nil {
		t.Fatalf(tplRunError, "create", err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "*_add_column.sql"))
	if len(matches) != 1 {
		t.Errorf("ожидался один файл *_add_column.sql, найдено %d", len(matches))
	}
	matches, _ = filepath.Glob(filepath.Join(dir, "*_fetch_data.go"))
	if len(matches) != 1 {
		t.Errorf("ожидался один файл *_fetch_data.go, найдено %d", len(matches))
	}
}

func TestRunErrors(t *testing.T) {
	const (
		errUnknown = "ожидалась ошибка неизвестной команды"
		errUpTo    = "ожидалась ошибка up-to без версии"
		errUpToNaN = "ожидалась ошибка нечисловой версии up-to"
		errDownTo  = "ожидалась ошибка down-to без версии"
		errDownNaN = "ожидалась ошибка нечисловой версии down-to"
		errCreate  = "ожидалась ошибка create без имени"
	)
	var (
		db  = testDatabase(t)
		dir = testThreeMigrationsDir(t)
		err error
	)

	if err = Run("no-such-command", db, dir); err == nil {
		t.Fatal(errUnknown)
	}
	if err = Run("up-to", db, dir); err == nil {
		t.Fatal(errUpTo)
	}
	if err = Run("up-to", db, dir, "abc"); err == nil {
		t.Fatal(errUpToNaN)
	}
	if err = Run("down-to", db, dir); err == nil {
		t.Fatal(errDownTo)
	}
	if err = Run("down-to", db, dir, "abc"); err == nil {
		t.Fatal(errDownNaN)
	}
	if err = Run("create", nil, dir); err == nil {
		t.Fatal(errCreate)
	}
}
