package main

import (
	"testing"
)

// Замена функции завершения процесса на подменную с записью кода.
func testTerminateSwap(t *testing.T) (code *int) {
	var old func(int)

	code = new(int)
	old = terminate
	terminate = func(c int) { *code = c }
	t.Helper()
	t.Cleanup(func() { terminate = old })

	return
}

func TestArgsDefaults(t *testing.T) {
	const tplWrongCommand = "ожидалась команда %q, получено %q"
	var (
		code *int
		cmd  string
		arg  *Args
		err  error
	)

	code = testTerminateSwap(t)
	if cmd, arg, err = args([]string{"status"}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "status" {
		t.Errorf(tplWrongCommand, "status", cmd)
	}
	if arg.Directory != "." || arg.Driver != "mysql" || arg.Dsn == "" {
		t.Errorf("ожидались значения по умолчанию: dir=%q drv=%q dsn=%q", arg.Directory, arg.Driver, arg.Dsn)
	}
	if *code != 0 {
		t.Errorf("не ожидалось завершение процесса с кодом %d", *code)
	}
}

func TestArgsFlags(t *testing.T) {
	var (
		cmd string
		arg *Args
		err error
	)

	if cmd, arg, err = args([]string{
		"--dir=/tmp/migration",
		"--drv=sqlite3",
		"--dsn=./foo.db",
		"up",
	}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "up" {
		t.Errorf("ожидалась команда up, получено %q", cmd)
	}
	if arg.Directory != "/tmp/migration" || arg.Driver != "sqlite3" || arg.Dsn != "./foo.db" {
		t.Errorf("неверные значения флагов: dir=%q drv=%q dsn=%q", arg.Directory, arg.Driver, arg.Dsn)
	}
}

func TestArgsEnVar(t *testing.T) {
	var (
		cmd string
		arg *Args
		err error
	)

	t.Setenv("GOOSE_DIR", "/tmp/envar_dir")
	t.Setenv("GOOSE_DRV", "sqlite3")
	t.Setenv("GOOSE_DSN", "./envar.db")
	if cmd, arg, err = args([]string{"version"}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "version" {
		t.Errorf("ожидалась команда version, получено %q", cmd)
	}
	if arg.Directory != "/tmp/envar_dir" || arg.Driver != "sqlite3" || arg.Dsn != "./envar.db" {
		t.Errorf("неверные значения переменных окружения: dir=%q drv=%q dsn=%q", arg.Directory, arg.Driver, arg.Dsn)
	}
}

func TestArgsUpDownVersion(t *testing.T) {
	var (
		cmd string
		arg *Args
		err error
	)

	if cmd, arg, err = args([]string{"up-to", "20180306100000"}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "up-to" || arg.UpDownVersion != "20180306100000" {
		t.Errorf("ожидалась команда up-to с версией 20180306100000, получено %q и %q", cmd, arg.UpDownVersion)
	}
	if cmd, arg, err = args([]string{"down-to", "20180305100000"}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "down-to" || arg.UpDownVersion != "20180305100000" {
		t.Errorf("ожидалась команда down-to с версией 20180305100000, получено %q и %q", cmd, arg.UpDownVersion)
	}
}

func TestArgsCreateDefaults(t *testing.T) {
	var (
		cmd string
		arg *Args
		err error
	)

	if cmd, arg, err = args([]string{"create"}); err != nil {
		t.Fatalf("не ожидалась ошибка разбора аргументов: %v", err)
	}
	if cmd != "create" || arg.CreateName != "new" || arg.CreateType != "sql" {
		t.Errorf("ожидалась команда create с именем new и типом sql, получено %q, %q, %q", cmd, arg.CreateName, arg.CreateType)
	}
}

func TestArgsHelp(t *testing.T) {
	var (
		code *int
		err  error
	)

	code = testTerminateSwap(t)
	if _, _, err = args([]string{"--help"}); err == nil {
		t.Fatal("ожидалась ошибка запроса справки")
	}
	if *code != 0 {
		t.Errorf("ожидался код завершения 0, получено %d", *code)
	}
}

func TestArgsNoCommand(t *testing.T) {
	var (
		code *int
		err  error
	)

	code = testTerminateSwap(t)
	if _, _, err = args([]string{}); err == nil {
		t.Fatal("ожидалась ошибка отсутствия команды")
	}
	if *code != 1 {
		t.Errorf("ожидался код завершения 1, получено %d", *code)
	}
}

func TestArgsBadFlag(t *testing.T) {
	var (
		code *int
		err  error
	)

	code = testTerminateSwap(t)
	if _, _, err = args([]string{"--no-such-flag", "up"}); err == nil {
		t.Fatal("ожидалась ошибка неизвестного флага")
	}
	if *code != 1 {
		t.Errorf("ожидался код завершения 1, получено %d", *code)
	}
}
