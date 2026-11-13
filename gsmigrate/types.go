package main

import (
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

const (
	appHelp = "Утилита применения миграций баз данных" // Описание утилиты в справке.
	// Описания флагов.
	flagHelpDirectory = "Каталог с файлами миграций. Значение по умолчанию — '.'. Переопределяет значение флага по умолчанию из переменной окружения с именем 'GOOSE_DIR'"
	flagHelpDriver    = "Драйвер базы данных. Поддерживаются [mysql, postgres, sqlite3, redshift, clickhouse, tidb]. Переопределяет значение флага по умолчанию из переменной окружения с именем 'GOOSE_DRV'"
	flagHelpDsn       = "Строка подключения к базе данных (DSN). Переопределяет значение флага по умолчанию из переменной окружения с именем 'GOOSE_DSN'"
	// Описания команд.
	cmdHelpUp      = "Применение к базе данных миграций до самой последней доступной версии"
	cmdHelpUpTo    = "Применение к базе данных миграций до указанной VERSION"
	cmdHelpDown    = "Откат версии на 1 назад"
	cmdHelpDownTo  = "Откат миграций до указанной VERSION"
	cmdHelpRedo    = "Повторное выполнение последней миграции"
	cmdHelpReset   = "Откат всех применённых миграций"
	cmdHelpStatus  = "Вывод статуса миграций для текущей базы данных"
	cmdHelpVersion = "Вывод текущей версии базы данных"
	cmdHelpCreate  = "Создание нового файла миграции со следующей версией"
	// Описания аргументов команд.
	argHelpVersion    = "Номер версии миграции"
	argHelpCreateName = "Имя файла новой миграции. Имя по умолчанию — 'new'"
	argHelpCreateType = "Тип файла новой миграции [sql|go]. Тип по умолчанию — 'sql'"
	// Значения по умолчанию.
	defaultDirectory = "."
	defaultDriver    = "mysql"
	defaultDsn       = `root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true`
	defaultName      = "new"
	defaultType      = "sql"
	// Имена переменных окружения.
	envDirectory = "GOOSE_DIR"
	envDriver    = "GOOSE_DRV"
	envDsn       = "GOOSE_DSN"
	// Имя приложения.
	appName = "gsmigrate"
)

var (
	terminate = func(code int) { os.Exit(code) }                                 // Завершение процесса с заданным кодом возврата (переопределяется в тестах).
	logFatalf = func(format string, args ...any) { log.Fatalf(format, args...) } // Вывод критической ошибки и завершение процесса (переопределяется в тестах).
)

// Args Аргументы командной строки утилиты.
type Args struct {
	Up            *kingpin.CmdClause // Команда применения всех миграций.
	UpTo          *kingpin.CmdClause // Команда применения миграций до версии.
	Down          *kingpin.CmdClause // Команда отката одной миграции.
	DownTo        *kingpin.CmdClause // Команда отката миграций до версии.
	UpDownVersion string             // Номер версии миграции для команд up-to и down-to.
	Redo          *kingpin.CmdClause // Команда повторного выполнения миграции.
	Reset         *kingpin.CmdClause // Команда отката всех миграций.
	Status        *kingpin.CmdClause // Команда вывода статуса миграций.
	Version       *kingpin.CmdClause // Команда вывода версии базы данных.
	Create        *kingpin.CmdClause // Команда создания файла миграции.
	CreateName    string             // Имя файла новой миграции.
	CreateType    string             // Тип файла новой миграции.
	Directory     string             // Каталог с файлами миграций.
	Driver        string             // Драйвер базы данных.
	Dsn           string             // Строка подключения к базе данных (DSN).
}
