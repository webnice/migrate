package goose

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

// Status Вывод статуса всех миграций.
func Status(db *sql.DB, dir string) (err error) {
	const (
		tplHeader    = "    Дата и время применения       -- Миграция"
		tplDelimiter = "    ========================================="
	)
	var (
		migrations Migrations
		mig        *Migration
	)

	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	// Создание таблицы версий при работе с чистой базой данных.
	if _, err = EnsureDBVersion(db); err != nil {
		return
	}
	log.Println(tplHeader)
	log.Println(tplDelimiter)
	for _, mig = range migrations {
		if err = printMigrationStatus(db, mig.Version, filepath.Base(mig.Source)); err != nil {
			return
		}
	}

	return
}

// Вывод статуса одной миграции.
func printMigrationStatus(db *sql.DB, version int64, script string) (err error) {
	const (
		tplQuery    = "SELECT tstamp, is_applied FROM goose_db_version WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1"
		tplPending  = "Ожидает"
		tplRow      = "    %-29s -- %s"
		tplDateTime = "02.01.2006 15:04:05 -0700 MST"
	)
	var (
		row       MigrationRecord
		appliedAt string
		q         string
		e         error
	)

	q = fmt.Sprintf(tplQuery, version)
	if e = db.QueryRow(q).Scan(&row.TStamp, &row.IsApplied); e != nil && !errors.Is(e, sql.ErrNoRows) {
		err = Errors().ScanMigrationStatus.Bind(version, e)
		return
	}
	if row.IsApplied {
		appliedAt = row.TStamp.In(time.Local).Format(tplDateTime)
	} else {
		appliedAt = tplPending
	}
	log.Printf(tplRow, appliedAt, script)

	return
}
