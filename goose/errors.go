package goose

import "github.com/webnice/dic"

// Все ошибки определены как константы.
const (
	cNoCurrentVersion         = "Миграция с текущей версией не найдена."
	cNoNextVersion            = "Следующая версия миграции не найдена."
	cUpToUsage                = "Команда up-to: требуется номер версии: goose [OPTIONS] DRIVER DBSTRING up-to VERSION."
	cCreateUsage              = "Команда create: требуется имя миграции: goose [OPTIONS] DRIVER DBSTRING create NAME [go|sql]."
	cDownToUsage              = "Команда down-to: требуется номер версии: goose [OPTIONS] DRIVER DBSTRING down-to VERSION."
	cNotMigrationFile         = "Не распознанный тип файла миграции."
	cNoSeparator              = "Разделитель не найден."
	cNotPositiveVersion       = "Номера миграций должны быть больше нуля."
	cNoAnnotations            = "Аннотации Up/Down не найдены, выражения не выполнялись."
	cDirectoryNotExists       = "Каталог %s не существует."
	cVersionNotNumber         = "Номер версии должен быть числом, получено %q."
	cUnknownCommand           = "%q: неизвестная команда."
	cMigrationConflict        = "Ошибка добавления миграции %q: версия совпадает с версией миграции %q."
	cMigrationFail            = "Не удалось выполнить миграцию %q, ошибка: %s."
	cGoMigrationNotRegistered = "Не удалось применить Go-миграцию %q: функции Go должны быть зарегистрированы и собраны в отдельный бинарный файл."
	cBeginTransaction         = "Ошибка открытия транзакции: %s."
	cMigrationFunc            = "Ошибка выполнения миграции %s, ошибка: %s."
	cInsertVersion            = "Ошибка записи версии миграции: %s."
	cParseVersion             = "Ошибка разбора номера версии: %s."
	cUnknownDialect           = "%q: неизвестный диалект."
	cFileAlreadyExists        = "Не удалось создать файл: %s уже существует."
	cCreateVersionTable       = "Ошибка создания таблицы версий миграций: %s."
	cInsertInitialVersion     = "Ошибка записи начальной версии миграций: %s."
	cScanVersionRecords       = "Ошибка чтения записей версий миграций: %s."
	cReadMigrationFile        = "Ошибка чтения файла миграции: %s."
	cExecStatement            = "Ошибка выполнения выражения: %s."
	cOpenMigrationFile        = "Ошибка открытия файла миграции %q, ошибка: %s."
	cMigrationNotFound        = "Миграция %d не найдена."
	cDuplicateVersion         = "Обнаружена дублирующаяся версия %d: %q и %q."
	cScanMigrationStatus      = "Ошибка чтения статуса миграции %d, ошибка: %s."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// NoCurrentVersion Миграция с текущей версией не найдена.
	NoCurrentVersion dic.IError

	// NoNextVersion Следующая версия миграции не найдена.
	NoNextVersion dic.IError

	// UpToUsage Команда up-to вызвана без номера версии.
	UpToUsage dic.IError

	// CreateUsage Команда create вызвана без имени миграции.
	CreateUsage dic.IError

	// DownToUsage Команда down-to вызвана без номера версии.
	DownToUsage dic.IError

	// NotMigrationFile Не распознанный тип файла миграции.
	NotMigrationFile dic.IError

	// NoSeparator Разделитель не найден.
	NoSeparator dic.IError

	// NotPositiveVersion Номера миграций должны быть больше нуля.
	NotPositiveVersion dic.IError

	// NoAnnotations Аннотации Up/Down не найдены, выражения не выполнялись.
	NoAnnotations dic.IError

	// DirectoryNotExists Каталог ... не существует.
	DirectoryNotExists dic.IError

	// VersionNotNumber Номер версии должен быть числом, получено ...
	VersionNotNumber dic.IError

	// UnknownCommand ...: неизвестная команда.
	UnknownCommand dic.IError

	// MigrationConflict Ошибка добавления миграции ...: версия совпадает с версией миграции ...
	MigrationConflict dic.IError

	// MigrationFail Не удалось выполнить миграцию ..., ошибка: ...
	MigrationFail dic.IError

	// GoMigrationNotRegistered Не удалось применить Go-миграцию ...: функции не зарегистрированы.
	GoMigrationNotRegistered dic.IError

	// BeginTransaction Ошибка открытия транзакции: ...
	BeginTransaction dic.IError

	// MigrationFunc Ошибка выполнения миграции ..., ошибка: ...
	MigrationFunc dic.IError

	// InsertVersion Ошибка записи версии миграции: ...
	InsertVersion dic.IError

	// ParseVersion Ошибка разбора номера версии: ...
	ParseVersion dic.IError

	// UnknownDialect ...: неизвестный диалект.
	UnknownDialect dic.IError

	// FileAlreadyExists Не удалось создать файл: ... уже существует.
	FileAlreadyExists dic.IError

	// CreateVersionTable Ошибка создания таблицы версий миграций: ...
	CreateVersionTable dic.IError

	// InsertInitialVersion Ошибка записи начальной версии миграций: ...
	InsertInitialVersion dic.IError

	// ScanVersionRecords Ошибка чтения записей версий миграций: ...
	ScanVersionRecords dic.IError

	// ReadMigrationFile Ошибка чтения файла миграции: ...
	ReadMigrationFile dic.IError

	// ExecStatement Ошибка выполнения выражения: ...
	ExecStatement dic.IError

	// OpenMigrationFile Ошибка открытия файла миграции ..., ошибка: ...
	OpenMigrationFile dic.IError

	// MigrationNotFound Миграция ... не найдена.
	MigrationNotFound dic.IError

	// DuplicateVersion Обнаружена дублирующаяся версия ...: ... и ...
	DuplicateVersion dic.IError

	// ScanMigrationStatus Ошибка чтения статуса миграции ..., ошибка: ...
	ScanMigrationStatus dic.IError
}

var errSingleton = &Error{
	Errors:                   dic.Error(),
	NoCurrentVersion:         dic.NewError(cNoCurrentVersion),
	NoNextVersion:            dic.NewError(cNoNextVersion),
	UpToUsage:                dic.NewError(cUpToUsage),
	CreateUsage:              dic.NewError(cCreateUsage),
	DownToUsage:              dic.NewError(cDownToUsage),
	NotMigrationFile:         dic.NewError(cNotMigrationFile),
	NoSeparator:              dic.NewError(cNoSeparator),
	NotPositiveVersion:       dic.NewError(cNotPositiveVersion),
	NoAnnotations:            dic.NewError(cNoAnnotations),
	DirectoryNotExists:       dic.NewError(cDirectoryNotExists, "каталог"),
	VersionNotNumber:         dic.NewError(cVersionNotNumber, "версия"),
	UnknownCommand:           dic.NewError(cUnknownCommand, "команда"),
	MigrationConflict:        dic.NewError(cMigrationConflict, "файл", "файл"),
	MigrationFail:            dic.NewError(cMigrationFail, "файл", "ошибка"),
	GoMigrationNotRegistered: dic.NewError(cGoMigrationNotRegistered, "файл"),
	BeginTransaction:         dic.NewError(cBeginTransaction, "ошибка"),
	MigrationFunc:            dic.NewError(cMigrationFunc, "файл", "ошибка"),
	InsertVersion:            dic.NewError(cInsertVersion, "ошибка"),
	ParseVersion:             dic.NewError(cParseVersion, "ошибка"),
	UnknownDialect:           dic.NewError(cUnknownDialect, "диалект"),
	FileAlreadyExists:        dic.NewError(cFileAlreadyExists, "файл"),
	CreateVersionTable:       dic.NewError(cCreateVersionTable, "ошибка"),
	InsertInitialVersion:     dic.NewError(cInsertInitialVersion, "ошибка"),
	ScanVersionRecords:       dic.NewError(cScanVersionRecords, "ошибка"),
	ReadMigrationFile:        dic.NewError(cReadMigrationFile, "ошибка"),
	ExecStatement:            dic.NewError(cExecStatement, "ошибка"),
	OpenMigrationFile:        dic.NewError(cOpenMigrationFile, "файл", "ошибка"),
	MigrationNotFound:        dic.NewError(cMigrationNotFound, "версия"),
	DuplicateVersion:         dic.NewError(cDuplicateVersion, "версия", "файл", "файл"),
	ScanMigrationStatus:      dic.NewError(cScanMigrationStatus, "версия", "ошибка"),
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
