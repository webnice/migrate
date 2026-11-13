package goose

import "database/sql"

// GetDialect Возвращает текущий диалект базы данных.
func GetDialect() SQLDialect { return dialect }

// SetDialect Устанавливает диалект базы данных по его имени.
func SetDialect(d string) (err error) {
	switch d {
	case "mysql":
		dialect = &MySQLDialect{}
	case "postgres":
		dialect = &PostgresDialect{}
	case "sqlite3":
		dialect = &Sqlite3Dialect{}
	case "redshift":
		dialect = &RedshiftDialect{}
	case "tidb":
		dialect = &TiDBDialect{}
	case "clickhouse":
		dialect = &ClickHouseDialect{}
	default:
		err = Errors().UnknownDialect.Bind(d)
	}

	return
}

// PostgresDialect Диалект PostgreSQL.
type PostgresDialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (pg PostgresDialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
              id serial NOT NULL,
              version_id bigint NOT NULL,
              is_applied boolean NOT NULL,
              tstamp timestamp NULL default now(),
              PRIMARY KEY(id)
            );`
}

// SQL-выражение добавления записи в таблицу версий.
func (pg PostgresDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, $2);"
}

// Запрос чтения лога версий миграций.
func (pg PostgresDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
}

// MySQLDialect Диалект MySQL.
type MySQLDialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (m MySQLDialect) createVersionTableSQL() string {
	return "CREATE TABLE `goose_db_version` (" +
		"id serial NOT NULL COMMENT 'Уникальный идентификатор записи.', " +
		"version_id BIGINT(20) NOT NULL COMMENT 'Номер версии базы данных.', " +
		"is_applied BOOLEAN NOT NULL COMMENT 'Истина - миграция применена. Ложь - миграция не применена.', " +
		"tstamp TIMESTAMP NULL DEFAULT now() COMMENT 'Дата и время миграции.', " +
		"PRIMARY KEY (id)" +
		") AUTO_INCREMENT = 1 COMMENT = 'Лог миграций';"
}

// SQL-выражение добавления записи в таблицу версий.
func (m MySQLDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

// Запрос чтения лога версий миграций.
func (m MySQLDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied FROM goose_db_version ORDER BY id DESC")
}

// Sqlite3Dialect Диалект sqlite3.
type Sqlite3Dialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (m Sqlite3Dialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
	          id INTEGER PRIMARY KEY AUTOINCREMENT,
	          version_id INTEGER NOT NULL,
              is_applied INTEGER NOT NULL,
              tstamp TIMESTAMP DEFAULT (datetime('now'))
            );`
}

// SQL-выражение добавления записи в таблицу версий.
func (m Sqlite3Dialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

// Запрос чтения лога версий миграций.
func (m Sqlite3Dialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
}

// RedshiftDialect Диалект Redshift.
type RedshiftDialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (rs RedshiftDialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
              id integer NOT NULL identity(1, 1),
              version_id bigint NOT NULL,
              is_applied boolean NOT NULL,
              tstamp timestamp NULL default sysdate,
              PRIMARY KEY(id)
            );`
}

// SQL-выражение добавления записи в таблицу версий.
func (rs RedshiftDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, $2);"
}

// Запрос чтения лога версий миграций.
func (rs RedshiftDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
}

// TiDBDialect Диалект TiDB.
type TiDBDialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (m TiDBDialect) createVersionTableSQL() string {
	return `CREATE TABLE goose_db_version (
	          id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE,
	          version_id bigint NOT NULL,
	          is_applied boolean NOT NULL,
	          tstamp timestamp NULL default now(),
	          PRIMARY KEY(id)
	        );`
}

// SQL-выражение добавления записи в таблицу версий.
func (m TiDBDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?);"
}

// Запрос чтения лога версий миграций.
func (m TiDBDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied from goose_db_version ORDER BY id DESC")
}

// ClickHouseDialect Диалект ClickHouse.
type ClickHouseDialect struct{}

// SQL-выражение создания таблицы goose_db_version.
func (ClickHouseDialect) createVersionTableSQL() string {
	return `
		CREATE TABLE goose_db_version (
		  version_id Int64,
		  is_applied UInt8,
		  date       Date     default today(),
		  tstamp     DateTime default now()
		) Engine = MergeTree(date, (date), 8192)
	`
}

// SQL-выражение добавления записи в таблицу версий.
func (ClickHouseDialect) insertVersionSQL() string {
	return "INSERT INTO goose_db_version (version_id, is_applied) VALUES (?, ?)"
}

// Запрос чтения лога версий миграций.
func (ClickHouseDialect) dbVersionQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT version_id, is_applied FROM goose_db_version ORDER BY version_id DESC, tstamp DESC")
}
