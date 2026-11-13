package goose

import (
	"errors"
	"testing"

	"github.com/webnice/dic"
)

func TestErrorsSingleton(t *testing.T) {
	const errSingletonBroken = "функция Errors(), функция повреждена"

	if Errors() != errSingleton {
		t.Fatal(errSingletonBroken)
	}
	if !Errors().NoCurrentVersion.Is(Errors().NoCurrentVersion.Bind()) {
		t.Error("ошибка NoCurrentVersion не эквивалентна справочнику")
	}
	if !Errors().NoNextVersion.Is(Errors().NoNextVersion.Bind()) {
		t.Error("ошибка NoNextVersion не эквивалентна справочнику")
	}
}

// TestErrors Тестирование всех ошибок пакета - корректность текста ошибки и эквивалентность справочнику.
//
//goland:noinspection DuplicatedCode
func TestErrors(t *testing.T) {
	const errBoom = "тестовая ошибка"
	var (
		base error
	)
	var items = []struct {
		name string
		got  error
		ref  dic.IError
		text string
	}{
		{"NoCurrentVersion", Errors().NoCurrentVersion.Bind(), Errors().NoCurrentVersion, cNoCurrentVersion},
		{"NoNextVersion", Errors().NoNextVersion.Bind(), Errors().NoNextVersion, cNoNextVersion},
		{"UpToUsage", Errors().UpToUsage.Bind(), Errors().UpToUsage, cUpToUsage},
		{"CreateUsage", Errors().CreateUsage.Bind(), Errors().CreateUsage, cCreateUsage},
		{"DownToUsage", Errors().DownToUsage.Bind(), Errors().DownToUsage, cDownToUsage},
		{"NotMigrationFile", Errors().NotMigrationFile.Bind(), Errors().NotMigrationFile, cNotMigrationFile},
		{"NoSeparator", Errors().NoSeparator.Bind(), Errors().NoSeparator, cNoSeparator},
		{"NotPositiveVersion", Errors().NotPositiveVersion.Bind(), Errors().NotPositiveVersion, cNotPositiveVersion},
		{"NoAnnotations", Errors().NoAnnotations.Bind(), Errors().NoAnnotations, cNoAnnotations},
		{"DirectoryNotExists", Errors().DirectoryNotExists.Bind("/tmp/migration"), Errors().DirectoryNotExists, "Каталог /tmp/migration не существует."},
		{"VersionNotNumber", Errors().VersionNotNumber.Bind("abc"), Errors().VersionNotNumber, `Номер версии должен быть числом, получено "abc".`},
		{"UnknownCommand", Errors().UnknownCommand.Bind("up"), Errors().UnknownCommand, `"up": неизвестная команда.`},
		{"MigrationConflict", Errors().MigrationConflict.Bind("a.go", "b.go"), Errors().MigrationConflict, `Ошибка добавления миграции "a.go": версия совпадает с версией миграции "b.go".`},
		{"GoMigrationNotRegistered", Errors().GoMigrationNotRegistered.Bind("a.go"), Errors().GoMigrationNotRegistered, `Не удалось применить Go-миграцию "a.go": функции Go должны быть зарегистрированы и собраны в отдельный бинарный файл.`},
		{"UnknownDialect", Errors().UnknownDialect.Bind("oracle"), Errors().UnknownDialect, `"oracle": неизвестный диалект.`},
		{"FileAlreadyExists", Errors().FileAlreadyExists.Bind("/tmp/x.sql"), Errors().FileAlreadyExists, "Не удалось создать файл: /tmp/x.sql уже существует."},
		{"MigrationNotFound", Errors().MigrationNotFound.Bind(20180305100000), Errors().MigrationNotFound, "Миграция 20180305100000 не найдена."},
		{"DuplicateVersion", Errors().DuplicateVersion.Bind(1, "a.sql", "b.sql"), Errors().DuplicateVersion, `Обнаружена дублирующаяся версия 1: "a.sql" и "b.sql".`},
	}

	for n := range items {
		if items[n].got.Error() != items[n].text {
			t.Errorf("функция ошибки %q, текст ошибки: %q, ожидался: %q", items[n].name, items[n].got.Error(), items[n].text)
		}
		if !items[n].ref.Is(items[n].got) {
			t.Errorf("функция ошибки %q, ошибка не эквивалентна справочнику", items[n].name)
		}
	}
	// Ошибки с обёрнутой исходной ошибкой: корректность текста и эквивалентность справочнику.
	base = errors.New(errBoom)
	var wrapped = []struct {
		name string
		got  error
		ref  dic.IError
	}{
		{"MigrationFail", Errors().MigrationFail.Bind("a.sql", base), Errors().MigrationFail},
		{"BeginTransaction", Errors().BeginTransaction.Bind(base), Errors().BeginTransaction},
		{"MigrationFunc", Errors().MigrationFunc.Bind("a.go", base), Errors().MigrationFunc},
		{"InsertVersion", Errors().InsertVersion.Bind(base), Errors().InsertVersion},
		{"ParseVersion", Errors().ParseVersion.Bind(base), Errors().ParseVersion},
		{"CreateVersionTable", Errors().CreateVersionTable.Bind(base), Errors().CreateVersionTable},
		{"InsertInitialVersion", Errors().InsertInitialVersion.Bind(base), Errors().InsertInitialVersion},
		{"ScanVersionRecords", Errors().ScanVersionRecords.Bind(base), Errors().ScanVersionRecords},
		{"ReadMigrationFile", Errors().ReadMigrationFile.Bind(base), Errors().ReadMigrationFile},
		{"ExecStatement", Errors().ExecStatement.Bind(base), Errors().ExecStatement},
		{"OpenMigrationFile", Errors().OpenMigrationFile.Bind("a.sql", base), Errors().OpenMigrationFile},
		{"ScanMigrationStatus", Errors().ScanMigrationStatus.Bind(1, base), Errors().ScanMigrationStatus},
	}

	for n := range wrapped {
		if !wrapped[n].ref.Is(wrapped[n].got) {
			t.Errorf("функция ошибки %q, ошибка не эквивалентна справочнику", wrapped[n].name)
		}
	}
}
