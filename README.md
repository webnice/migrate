# migrate

---

[![PkgGoDev](https://pkg.go.dev/badge/github.com/webnice/migrate/v2)](https://pkg.go.dev/github.com/webnice/migrate/v2)
[![Coverage Status](https://coveralls.io/repos/github/webnice/migrate/badge.svg?branch=v2)](https://coveralls.io/github/webnice/migrate?branch=v2)

### migrate

Инструмент миграции баз данных: управление схемой базы данных через инкрементные
изменения на SQL или функции Go.

---

## Описание

`migrate` — это инструмент управления схемой базы данных, состоящий из
библиотеки `goose` и готовой утилиты командной строки `gsmigrate`. Изменения
схемы оформляются как миграции — отдельные файлы с номером версии, которые
применяются последовательно и могут быть откачены. Состояние базы данных
отслеживается в служебной таблице `goose_db_version`.

### Особенности

- Применение и откат миграций: по одной, все до последней версии, до указанной версии;
- SQL-миграции с аннотациями `-- +goose Up` и `-- +goose Down`;
- Поддержка сложных выражений PL/pgSQL через аннотации `StatementBegin` и `StatementEnd`;
- Выполнение миграций внутри транзакции с возможностью отключения (`NO TRANSACTION`);
- Go-миграции: регистрация функций через `goose.AddMigration` в собственном бинарном файле;
- Команды `status` и `version` для контроля состояния базы данных;
- Значения флагов по умолчанию из переменных окружения `GOOSE_DIR`, `GOOSE_DRV`, `GOOSE_DSN`;
- Создание файлов новых миграций командой `create` без подключения к базе данных;
- Цели `m-[driver]-[command]` в Makefile проекта для работы с миграциями;
- Покрытие кода тестами 100% (пакет `goose`).

## Требования

- Go 1.24 или новее.

## Установка

```bash
go install github.com/webnice/migrate/v2/gsmigrate@latest
```

Подключение библиотеки:

```go
import (
	"github.com/webnice/migrate/v2/goose"
)
```

#### Зависимости

- `github.com/alecthomas/kingpin/v2` — разбор аргументов командной строки;
- `github.com/ClickHouse/clickhouse-go` — драйвер ClickHouse;
- `github.com/go-sql-driver/mysql` — драйвер MySQL;
- `github.com/lib/pq` — драйвер PostgreSQL, CockroachDB, Redshift;
- `github.com/mattn/go-sqlite3` — драйвер SQLite3;
- `github.com/ziutek/mymysql` — драйвер App Engine CloudSQL.

## Поддерживаемые базы данных

| Диалект     | Драйвер    | Пример DSN |
|-------------|------------|------------|
| mysql       | mysql      | `user:password@/dbname?parseTime=true` |
| postgres    | postgres   | `user=postgres dbname=postgres sslmode=disable` |
| cockroach   | postgres   | `user=root dbname=default sslmode=disable` |
| sqlite3     | sqlite3    | `./foo.db` |
| redshift    | postgres   | `postgres://user:password@host:5439/db` |
| clickhouse  | clickhouse | `tcp://localhost:9000?username=default&database=test` |
| tidb        | mysql      | `user:password@/dbname?parseTime=true` |

## Использование

```
usage: gsmigrate [<flags>] <command> [<args> ...]

Утилита применения миграций баз данных

Flags:
  --dir="."      Каталог с файлами миграций. Значение по умолчанию — '.'.
                 Переопределяется значение флага по умолчанию из переменной
                 окружения с именем 'GOOSE_DIR'
  --drv="mysql"  Драйвер базы данных. Поддерживаются [mysql, postgres, sqlite3,
                 redshift, clickhouse, tidb]. Переопределяет значение флага по
                 умолчанию из переменной окружения с именем 'GOOSE_DRV'
  --dsn="root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true"
                 Строка подключения к базе данных (DSN). Переопределяет
                 значение флага по умолчанию из переменной окружения с именем
                 'GOOSE_DSN'

Commands:
  up             Применение к базе данных миграций до самой последней
                 доступной версии
  up-to [<VERSION>]
                 Применение к базе данных миграций до указанной VERSION
  down           Откат версии на 1 назад
  down-to [<VERSION>]
                 Откат миграций до указанной VERSION
  redo           Повторное выполнение последней миграции
  reset          Откат всех применённых миграций
  status         Вывод статуса миграций для текущей базы данных
  version        Вывод текущей версии базы данных
  create [<NAME>] [<TYPE>]
                 Создание нового файла миграции со следующей версией
```

Полную справку по командам выводит `gsmigrate --help`, справку по отдельной
команде — `gsmigrate help <command>`.

## Примеры

Нижеприведённая сессия выполняется с базой данных sqlite3.

Создание миграции и просмотр статуса:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db create init sql
Создан новый файл: 20261113071348_init.sql

$ gsmigrate --drv=sqlite3 --dsn=./foo.db status
    Дата и время применения       -- Миграция
    =========================================
    Ожидает                       -- 20261113071348_init.sql
```

Применение всех миграций и вывод версии:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db up
OK 20261113071348_init.sql

$ gsmigrate --drv=sqlite3 --dsn=./foo.db version
goose: версия 20261113071348
```

Применение до указанной версии:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db up-to 20180306100000
OK 20180305100000_begin.sql
OK 20180306100000_next.sql
goose: нет миграций для выполнения. Текущая версия: 20180306100000
```

Откат одной миграции и откат до версии:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db down
OK 20180307100000_and_again.sql

$ gsmigrate --drv=sqlite3 --dsn=./foo.db down-to 20180305100000
OK 20180307100000_and_again.sql
OK 20180306100000_next.sql
goose: нет миграций для выполнения. Текущая версия: 20180305100000
```

Повторное выполнение последней миграции:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db redo
OK 20180307100000_and_again.sql
OK 20180307100000_and_again.sql
```

Полный откат всех миграций:

```bash
$ gsmigrate --drv=sqlite3 --dsn=./foo.db reset
OK 20180307100000_and_again.sql
OK 20180306100000_next.sql
OK 20180305100000_begin.sql
```

## SQL-миграции

Файл SQL-миграции содержит секции, размеченные аннотациями goose в комментариях:

```sql
-- +goose Up
CREATE TABLE post (
    id int NOT NULL,
    title text,
    body text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

Выражения после `-- +goose Up` выполняются при применении миграции, после
`-- +goose Down` — при откате. По умолчанию все миграции выполняются внутри
транзакции. Выражения, которые не могут выполняться внутри транзакции
(например `CREATE DATABASE`), размещаются в файле с аннотацией
`-- +goose NO TRANSACTION` — тогда обе секции файла выполняются без транзакции.

По умолчанию выражения разделяются точкой с запятой и должны заканчиваться
ею. Сложные выражения (PL/pgSQL), содержащие точки с запятой внутри себя,
обрамляются аннотациями `-- +goose StatementBegin` и `-- +goose StatementEnd`:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION histories_partition_creation(DATE, DATE)
returns void AS $$
DECLARE
  create_query text;
BEGIN
  FOR create_query IN SELECT
      'CREATE TABLE IF NOT EXISTS histories_'
      || TO_CHAR(d, 'YYYY_MM')
      || ' ( CHECK( created_at >= timestamp '''
      || TO_CHAR(d, 'YYYY-MM-DD 00:00:00')
      || ''' AND created_at < timestamp '''
      || TO_CHAR(d + INTERVAL '1 month', 'YYYY-MM-DD 00:00:00')
      || ''' ) ) inherits (histories);'
    FROM generate_series($1, $2, '1 month') AS d
  LOOP
    EXECUTE create_query;
  END LOOP;
END;
$$
language plpgsql;
-- +goose StatementEnd
```

Имя файла миграции должно начинаться с номера версии и содержать разделитель:
`XXX_descriptivename.sql`, где `XXX` — номер версии. Новые файлы создаются
командой `create`, версия формируется из текущего времени в формате UTC.

## Go-миграции

Go-миграции выполняются собственным бинарным файлом, собранным с пакетом
`goose`. Для этого:

1. Создайте собственный бинарный файл;
2. Импортируйте `github.com/webnice/migrate/v2/goose`;
3. Зарегистрируйте функции миграций;
4. Вызовите команду goose, например `goose.Up(db, dir)`.

Пример Go-миграции:

```go
package main

import (
	"database/sql"

	"github.com/webnice/migrate/v2/goose"
)

func init() {
	goose.AddMigration(Up20180307100000, Down20180307100000)
}

// Up20180307100000 Миграция применена.
func Up20180307100000(tx *sql.Tx) (err error) {
	// Этот код выполняется при применении миграции.
	return
}

// Down20180307100000 Миграция откачена.
func Down20180307100000(tx *sql.Tx) (err error) {
	// Этот код выполняется при откате миграции.
	return
}
```

Файл Go-миграции создаётся командой `gsmigrate create <NAME> go`. Функция
миграции получает открытую транзакцию; запись о версии вносится в лог
автоматически после успешного выполнения функции.

## Интеграция с Makefile проекта

Makefile проекта поддерживает цели вида `m-[driver]-[command]`. Драйверы и
параметры подключения описываются в файле `.env.mk` проекта:

```make
MIGRATION_DRIVER = mysql clickhouse

MIGRATION_DIR_MYSQL = migration/mysql
MIGRATION_DRV_MYSQL = mysql
MIGRATION_DSN_MYSQL = user:password@/dbname?parseTime=true

MIGRATION_DIR_CLICKHOUSE = migration/clickhouse
MIGRATION_DRV_CLICKHOUSE = clickhouse
MIGRATION_DSN_CLICKHOUSE = tcp://localhost:9000?username=default&database=test
```

После этого доступны цели:

```bash
make m-mysql-up        # применение всех миграций MySQL
make m-mysql-down      # откат одной миграции MySQL
make m-mysql-create    # создание новой миграции MySQL
make m-mysql-status    # статус миграций MySQL
make m-mysql-redo      # повтор последней миграции MySQL
make m-mysql-version   # версия базы данных MySQL
make m-clickhouse-up   # и так далее для каждого драйвера
```

## Тестирование

```bash
# Стандартный прогон с покрытием.
go test -cover ./...

# Прогон через Makefile (агрегированный профиль покрытия в coverage.log).
make test
```

## Лицензия

© WEB NICE, распространяется под [лицензией MIT](LICENSE).
