package main

import "github.com/alecthomas/kingpin/v2"

// Разбор аргументов командной строки.
func args(argv []string) (cmd string, ret *Args, err error) {
	const (
		flagNameDirectory = "dir"
		flagNameDriver    = "drv"
		flagNameDsn       = "dsn"
		cmdNameUp         = "up"
		cmdNameUpTo       = "up-to"
		cmdNameDown       = "down"
		cmdNameDownTo     = "down-to"
		cmdNameRedo       = "redo"
		cmdNameReset      = "reset"
		cmdNameStatus     = "status"
		cmdNameVersion    = "version"
		cmdNameCreate     = "create"
		argNameVersion    = "VERSION"
		argNameCreateName = "NAME"
		argNameCreateType = "TYPE"
	)
	var (
		app  *kingpin.Application
		help bool
		arg  string
	)

	ret = new(Args)
	// Определение запроса справки до разбора аргументов.
	for _, arg = range argv {
		if arg == "--help" || arg == "-h" {
			help = true
		}
	}
	app = kingpin.New(appName, appHelp)
	// Завершение процесса выполняется только из точки входа.
	app.Terminate(terminate)
	// Глобальные флаги.
	app.Flag(flagNameDirectory, flagHelpDirectory).
		Envar(envDirectory).
		Default(defaultDirectory).
		StringVar(&ret.Directory)
	app.Flag(flagNameDriver, flagHelpDriver).
		Envar(envDriver).
		Default(defaultDriver).
		StringVar(&ret.Driver)
	app.Flag(flagNameDsn, flagHelpDsn).
		Envar(envDsn).
		Default(defaultDsn).
		StringVar(&ret.Dsn)
	// Команды с аргументами.
	ret.Up = app.Command(cmdNameUp, cmdHelpUp)
	ret.UpTo = app.Command(cmdNameUpTo, cmdHelpUpTo)
	ret.UpTo.Arg(argNameVersion, argHelpVersion).
		Default("").
		StringVar(&ret.UpDownVersion)
	ret.Down = app.Command(cmdNameDown, cmdHelpDown)
	ret.DownTo = app.Command(cmdNameDownTo, cmdHelpDownTo)
	ret.DownTo.Arg(argNameVersion, argHelpVersion).
		Default("").
		StringVar(&ret.UpDownVersion)
	ret.Redo = app.Command(cmdNameRedo, cmdHelpRedo)
	ret.Reset = app.Command(cmdNameReset, cmdHelpReset)
	ret.Status = app.Command(cmdNameStatus, cmdHelpStatus)
	ret.Version = app.Command(cmdNameVersion, cmdHelpVersion)
	ret.Create = app.Command(cmdNameCreate, cmdHelpCreate)
	ret.Create.Arg(argNameCreateName, argHelpCreateName).
		Default(defaultName).
		StringVar(&ret.CreateName)
	ret.Create.Arg(argNameCreateType, argHelpCreateType).
		Default(defaultType).
		StringVar(&ret.CreateType)
	if cmd, err = app.Parse(argv); err != nil {
		if help {
			// Справка выведена kingpin — завершение без ошибки.
			terminate(0)
			return
		}
		terminate(1)
	}

	return
}
