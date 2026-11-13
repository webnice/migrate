/*

Утилита gsmigrate — командный интерфейс применения миграций базы данных.

Утилита построена на библиотеке github.com/webnice/migrate/goose и
предоставляет доступ ко всем её командам: up, up-to, down, down-to, redo,
reset, status, version и create.

Использование:

	gsmigrate [<флаги>] <команда> [<аргументы>]

Флаги:

	--dir  Каталог с файлами миграций, по умолчанию '.'.
	--drv  Драйвер базы данных: mysql, postgres, sqlite3, redshift,
	       clickhouse или tidb.
	--dsn  Строка подключения к базе данных (DSN).

Значения флагов по умолчанию переопределяются переменными окружения
GOOSE_DIR, GOOSE_DRV и GOOSE_DSN.

Команда create создаёт файл новой миграции и выполняется без подключения
к базе данных:

	gsmigrate --dir=./migration create add_some_column sql

Остальные команды применяют или откатывают миграции в базе данных:

	gsmigrate --drv=sqlite3 --dsn=./foo.db up

*/

package main
