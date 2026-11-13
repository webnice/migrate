package main

import "github.com/webnice/dic"

// Все ошибки определены как константы.
const (
	cUnknownDriver   = "%q: драйвер не поддерживается."
	cConnectDatabase = "Ошибка подключения к базе данных: %s."
	cRunCommand      = "Ошибка выполнения команды %q, ошибка: %s."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// UnknownDriver ...: драйвер не поддерживается.
	UnknownDriver dic.IError

	// ConnectDatabase Ошибка подключения к базе данных: ...
	ConnectDatabase dic.IError

	// RunCommand Ошибка выполнения команды ..., ошибка: ...
	RunCommand dic.IError
}

var errSingleton = &Error{
	Errors:          dic.Error(),
	UnknownDriver:   dic.NewError(cUnknownDriver, "драйвер"),
	ConnectDatabase: dic.NewError(cConnectDatabase, "ошибка"),
	RunCommand:      dic.NewError(cRunCommand, "команда", "ошибка"),
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
