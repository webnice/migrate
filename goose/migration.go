package goose

import (
	"database/sql"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// String Возвращает путь к файлу миграции.
func (m *Migration) String() string { return m.Source }

// Up Применение миграции.
func (m *Migration) Up(db *sql.DB) (err error) {
	const tplOK = "OK %s"

	if err = m.run(db, true); err != nil {
		return
	}
	log.Printf(tplOK, filepath.Base(m.Source))

	return
}

// Down Откат миграции.
func (m *Migration) Down(db *sql.DB) (err error) {
	const tplOK = "OK %s"

	if err = m.run(db, false); err != nil {
		return
	}
	log.Printf(tplOK, filepath.Base(m.Source))

	return
}

// Выполнение миграции в указанном направлении.
func (m *Migration) run(db *sql.DB, direction bool) (err error) {
	var (
		tx *sql.Tx
		fn func(*sql.Tx) error
	)

	switch filepath.Ext(m.Source) {
	case extensionSQL:
		if err = runSQLMigration(db, m.Source, m.Version, direction); err != nil {
			err = Errors().MigrationFail.Bind(m.Source, err)
			return
		}
	case extensionGo:
		if !m.Registered {
			err = Errors().GoMigrationNotRegistered.Bind(m.Source)
			return
		}
		if tx, err = db.Begin(); err != nil {
			err = Errors().BeginTransaction.Bind(err)
			return
		}
		fn = m.UpFn
		if !direction {
			fn = m.DownFn
		}
		if fn != nil {
			if err = fn(tx); err != nil {
				_ = tx.Rollback()
				err = Errors().MigrationFunc.Bind(filepath.Base(m.Source), err)
				return
			}
		}
		if _, err = tx.Exec(GetDialect().insertVersionSQL(), m.Version, direction); err != nil {
			_ = tx.Rollback()
			err = Errors().InsertVersion.Bind(err)
			return
		}
		err = tx.Commit()
	}

	return
}

// NumericComponent Извлекает номер версии из имени файла миграции вида
// XXX_descriptivename.ext, где XXX — номер версии, ext — тип миграции.
func NumericComponent(name string) (ret int64, err error) {
	var (
		base string
		ext  string
		idx  int
	)

	base = filepath.Base(name)
	if ext = filepath.Ext(base); ext != extensionSQL && ext != extensionGo {
		err = Errors().NotMigrationFile.Bind()
		return
	}
	if idx = strings.Index(base, "_"); idx < 0 {
		err = Errors().NoSeparator.Bind()
		return
	}
	if ret, err = strconv.ParseInt(base[:idx], 10, 64); err != nil {
		err = Errors().ParseVersion.Bind(err)
		return
	}
	if ret <= 0 {
		err = Errors().NotPositiveVersion.Bind()
	}

	return
}
