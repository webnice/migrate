package main

import (
//	"database/sql"
//	"os/exec"
//	"strings"
//	"testing"
//
//	_ "github.com/kshvakov/clickhouse"
//	"github.com/stretchr/testify/assert"
)

//func TestClickHouseDialect(t *testing.T) {
//	out, err := exec.Command("go", strings.Fields("go build -i -o gsmigrate github.com/webnice/migrate/gsmigrate")...).CombinedOutput()
//	if !assert.NoError(t, err, string(out)) {
//		return
//	}
//	if connect, err := sql.Open("clickhouse", "native://127.0.0.1:9000?debug=true"); assert.NoError(t, err) && assert.NoError(t, connect.Ping()) {
//		if _, err := connect.Exec("DROP DATABASE IF EXISTS goose_test"); !assert.NoError(t, err) {
//			return
//		}
//		if _, err := connect.Exec("CREATE DATABASE IF NOT EXISTS goose_test"); !assert.NoError(t, err) {
//			return
//		}
//
//		row := connect.QueryRow("SELECT COUNT() FROM system.tables WHERE database = 'goose_test' AND name = 'goose_db_version'")
//		var count int
//		if err := row.Scan(&count); !assert.Equal(t, sql.ErrNoRows, err) {
//			return
//		}
//		{
//			out, err := exec.Command("./gsmigrate", strings.Fields(`clickhouse native://127.0.0.1:9000?debug=true&database=goose_test status`)...).CombinedOutput()
//			if !assert.NoError(t, err, string(out)) {
//				return
//			}
//
//			var count int
//			row := connect.QueryRow("SELECT COUNT() FROM system.tables WHERE database = 'goose_test' AND name = 'goose_db_version'")
//			if err := row.Scan(&count); !assert.NoError(t, err) {
//				return
//			}
//			if !assert.Equal(t, int(1), count) {
//				return
//			}
//		}
//
//		var (
//			countOfColumns    int
//			countOfMigrations int
//		)
//
//		if err := connect.QueryRow("SELECT COUNT() FROM goose_test.goose_db_version").Scan(&countOfMigrations); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(1), countOfMigrations) {
//			return
//		}
//		{
//			out, err := exec.Command("./gsmigrate", strings.Fields(`-dir examples/sql-migrations/clickhouse clickhouse native://127.0.0.1:9000?debug=true&database=goose_test up`)...).CombinedOutput()
//			if !assert.NoError(t, err, string(out)) {
//				return
//			}
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM system.columns WHERE database = 'goose_test' AND table = 'events'").Scan(&countOfColumns); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(4), countOfColumns) {
//			return
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM system.columns WHERE database = 'goose_test' AND table = 'events' AND name = 'CreatedAt'").Scan(&countOfColumns); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(1), countOfColumns) {
//			return
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM goose_test.goose_db_version").Scan(&countOfMigrations); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(3), countOfMigrations) {
//			return
//		}
//		{
//			out, err := exec.Command("./gsmigrate", strings.Fields(`-dir examples/sql-migrations/clickhouse clickhouse native://127.0.0.1:9000?debug=true&database=goose_test down`)...).CombinedOutput()
//			if !assert.NoError(t, err, string(out)) {
//				return
//			}
//		}
//		if err := connect.QueryRow("SELECT COUNT() FROM system.columns WHERE database = 'goose_test' AND table = 'events'").Scan(&countOfColumns); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(3), countOfColumns) {
//			return
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM system.columns WHERE database = 'goose_test' AND table = 'events' AND name = 'CreatedAt'").Scan(&countOfColumns); !assert.Equal(t, sql.ErrNoRows, err) {
//			return
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM goose_test.goose_db_version").Scan(&countOfMigrations); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(4), countOfMigrations) {
//			return
//		}
//
//		if err := connect.QueryRow("SELECT COUNT() FROM goose_test.goose_db_version WHERE is_applied = 1").Scan(&countOfMigrations); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(3), countOfMigrations) {
//			return
//		}
//		if err := connect.QueryRow("SELECT COUNT() FROM goose_test.goose_db_version WHERE is_applied = 0").Scan(&countOfMigrations); !assert.NoError(t, err) {
//			return
//		}
//		if !assert.Equal(t, int(1), countOfMigrations) {
//			return
//		}
//	}
//}
