package goose

import (
	"database/sql"
	"fmt"
	"sort"
)

// Reset Откат всех применённых миграций.
func Reset(db *sql.DB, dir string) (err error) {
	var (
		migrations Migrations
		statuses   map[int64]bool
		mig        *Migration
		end        bool
	)

	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	if statuses, err = dbMigrationsStatus(db); err != nil {
		return
	}
	sort.Sort(sort.Reverse(migrations))
	for _, mig = range migrations {
		if end {
			break
		}
		if !statuses[mig.Version] {
			// Миграция не применялась — пропуск.
			continue
		}
		if err = mig.Down(db); err != nil {
			end = true
		}
	}

	return
}

// Возвращает состояние применённости всех версий миграций.
// При отсутствии таблицы версий создаёт и инициализирует её.
func dbMigrationsStatus(db *sql.DB) (ret map[int64]bool, err error) {
	const tplScan = "ошибка чтения записей версий миграций: %w"
	var (
		rows *sql.Rows
		row  MigrationRecord
		ok   bool
	)

	if rows, err = GetDialect().dbVersionQuery(db); err != nil {
		// Таблица версий отсутствует — создание и инициализация.
		err = createVersionTable(db)
		return
	}
	defer func() { _ = rows.Close() }()
	ret = make(map[int64]bool)
	for rows.Next() {
		if err = rows.Scan(&row.VersionID, &row.IsApplied); err != nil {
			err = fmt.Errorf(tplScan, err)
			return
		}
		if _, ok = ret[row.VersionID]; ok {
			// Учитывается самая свежая запись версии.
			continue
		}
		ret[row.VersionID] = row.IsApplied
	}

	return
}
