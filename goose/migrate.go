package goose

import (
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// Len Возвращает количество миграций в списке.
func (ms Migrations) Len() int { return len(ms) }

// Swap Меняет местами две миграции списка.
func (ms Migrations) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }

// Less Сравнивает номера версий двух миграций списка.
func (ms Migrations) Less(i, j int) bool { return ms[i].Version < ms[j].Version }

// Current Возвращает миграцию с указанной версией.
func (ms Migrations) Current(current int64) (ret *Migration, err error) {
	for i, migration := range ms {
		if migration.Version == current {
			ret = ms[i]
			return
		}
	}
	err = Errors().NoCurrentVersion.Bind()

	return
}

// Next Возвращает миграцию следующую за указанной версией.
func (ms Migrations) Next(current int64) (ret *Migration, err error) {
	for i, migration := range ms {
		if migration.Version > current {
			ret = ms[i]
			return
		}
	}
	err = Errors().NoNextVersion.Bind()

	return
}

// Previous Возвращает миграцию предыдущую перед указанной версией.
func (ms Migrations) Previous(current int64) (ret *Migration, err error) {
	for _, m := range slices.Backward(ms) {
		if m.Version < current {
			ret = m
			return
		}
	}
	err = Errors().NoNextVersion.Bind()

	return
}

// Last Возвращает последнюю миграцию списка.
func (ms Migrations) Last() (ret *Migration, err error) {
	if len(ms) == 0 {
		err = Errors().NoNextVersion.Bind()
		return
	}
	ret = ms[len(ms)-1]

	return
}

// String Возвращает список миграций, по одному файлу в строке.
func (ms Migrations) String() (ret string) {
	var (
		sb  strings.Builder
		mig *Migration
	)

	for _, mig = range ms {
		sb.WriteString(mig.String())
		sb.WriteString("\n")
	}
	ret = sb.String()

	return
}

// CollectMigrations Возвращает все корректные миграции из каталога и реестра Go-функций
// в диапазоне версий от current до target.
func CollectMigrations(dirpath string, current, target int64) (migrations Migrations, err error) {
	var (
		sqlMigrationFiles []string
		goMigrationFiles  []string
		file              string
		v                 int64
		mig               *Migration
		ok                bool
	)

	if _, err = os.Stat(dirpath); os.IsNotExist(err) {
		err = Errors().DirectoryNotExists.Bind(dirpath)
		return
	}
	// SQL-файлы миграций.
	sqlMigrationFiles = collectMigrationFiles(dirpath, extensionSQL)
	for _, file = range sqlMigrationFiles {
		if v, err = NumericComponent(file); err != nil {
			return
		}
		if versionFilter(v, current, target) {
			mig = &Migration{Version: v, Next: -1, Previous: -1, Source: file}
			migrations = append(migrations, mig)
		}
	}
	// Go-миграции, зарегистрированные через AddMigration().
	for _, mig = range registeredGoMigrations {
		if v, err = NumericComponent(mig.Source); err != nil {
			return
		}
		if versionFilter(v, current, target) {
			migrations = append(migrations, mig)
		}
	}
	// Go-файлы миграций без зарегистрированных функций.
	goMigrationFiles = collectMigrationFiles(dirpath, extensionGo)
	for _, file = range goMigrationFiles {
		if v, err = NumericComponent(file); err != nil {
			// Пропуск файлов без префикса версии в имени.
			continue
		}
		// Пропуск миграций, уже добавленных в реестр через AddMigration().
		if _, ok = registeredGoMigrations[v]; ok {
			continue
		}
		if versionFilter(v, current, target) {
			mig = &Migration{Version: v, Next: -1, Previous: -1, Source: file, Registered: false}
			migrations = append(migrations, mig)
		}
	}
	// Сортировка списка и обнаружение дублирующихся версий.
	if migrations, err = sortAndConnectMigrations(migrations); err != nil {
		return
	}

	return
}

// Поиск файлов миграций с указанным расширением в каталоге.
// Ошибка разбора шаблона невозможна: шаблон формируется из констант пакета.
func collectMigrationFiles(dirpath, extension string) (files []string) {
	files, _ = filepath.Glob(dirpath + "/**" + extension)

	return
}

// Сортировка списка миграций по версии и связывание миграций
// указателями следующей и предыдущей версии.
func sortAndConnectMigrations(migrations Migrations) (ret Migrations, err error) {
	var (
		seen map[int64]string
		src  string
		ok   bool
		i    int
		prev int64
		mig  *Migration
	)

	seen = make(map[int64]string)
	// Поиск миграций с одинаковым номером версии.
	for _, mig = range migrations {
		if src, ok = seen[mig.Version]; ok {
			err = Errors().DuplicateVersion.Bind(mig.Version, src, mig.Source)
			return
		}
		seen[mig.Version] = mig.Source
	}
	sort.Sort(migrations)
	// Заполнение связей следующей и предыдущей версии.
	for i, mig = range migrations {
		prev = -1
		if i > 0 {
			prev = migrations[i-1].Version
			migrations[i-1].Next = mig.Version
		}
		mig.Previous = prev
	}
	ret = migrations

	return
}

// Проверка попадания версии миграции в диапазон между current и target.
func versionFilter(v, current, target int64) (ret bool) {
	switch {
	case target > current:
		ret = v > current && v <= target
	case target < current:
		ret = v <= current && v > target
	}

	return
}

// EnsureDBVersion Возвращает текущую версию базы данных.
// При отсутствии таблицы версий создаёт и инициализирует её.
func EnsureDBVersion(db *sql.DB) (version int64, err error) {
	var (
		rows   *sql.Rows
		row    *MigrationRecord
		toSkip []int64
		skip   bool
		v      int64
	)

	// Чтение лога версий миграций.
	if rows, err = GetDialect().dbVersionQuery(db); err != nil {
		// Таблица версий отсутствует — создание и инициализация.
		err = createVersionTable(db)
		return
	}
	defer func() { _ = rows.Close() }()
	toSkip = make([]int64, 0)
	for rows.Next() {
		row = new(MigrationRecord)
		if err = rows.Scan(&row.VersionID, &row.IsApplied); err != nil {
			err = Errors().ScanVersionRecords.Bind(err)
			return
		}
		// Проверка наличия версии в списке пропускаемых.
		skip = false
		for _, v = range toSkip {
			if v == row.VersionID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		// Первая применённая версия — текущая.
		if row.IsApplied {
			version = row.VersionID
			return
		}
		// Последняя запись версии не применена — версия пропускается.
		toSkip = append(toSkip, row.VersionID)
	}
	err = Errors().NoNextVersion.Bind()

	return
}

// Создание таблицы goose_db_version
// и запись в неё начальной нулевой версии.
func createVersionTable(db *sql.DB) (err error) {
	var (
		txn     *sql.Tx
		d       SQLDialect
		version int
		applied bool
	)

	if txn, err = db.Begin(); err != nil {
		return
	}
	d = GetDialect()
	if _, err = txn.Exec(d.createVersionTableSQL()); err != nil {
		_ = txn.Rollback()
		err = Errors().CreateVersionTable.Bind(err)
		return
	}
	version, applied = 0, true
	if _, err = txn.Exec(d.insertVersionSQL(), version, applied); err != nil {
		_ = txn.Rollback()
		err = Errors().InsertInitialVersion.Bind(err)
		return
	}
	err = txn.Commit()

	return
}

// GetDBVersion Возвращает текущую версию базы данных.
// При ошибке чтения версия устанавливается в -1.
func GetDBVersion(db *sql.DB) (version int64, err error) {
	if version, err = EnsureDBVersion(db); err != nil {
		version = -1
	}

	return
}
