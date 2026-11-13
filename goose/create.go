package goose

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// CreateWithTemplate Создание нового файла миграции по заданному шаблону.
// Аргумент db зарезервирован для совместимости интерфейса и не используется.
func CreateWithTemplate(_ *sql.DB, dir string, migrationTemplate *template.Template, name, migrationType string) (err error) {
	const (
		versionLayout   = "20060102150405"
		tplFileName     = "%s_%s.%s"
		tplVersion      = "%014d"
		migrationTypeGo = "go"
		tplCreatedFile  = "Создан новый файл: %s"
	)
	var (
		migrations Migrations
		last       *Migration
		version    string
		filename   string
		fPath      string
		path       string
		tmpl       *template.Template
	)

	if migrations, err = CollectMigrations(dir, minVersion, MaxVersion); err != nil {
		return
	}
	// Формирование версии новой миграции из текущего времени в формате UTC.
	version = time.Now().In(time.UTC).Format(versionLayout)
	if last, err = migrations.Last(); err == nil {
		// Защита от повторной генерации версии в ту же секунду.
		if fmt.Sprintf(tplVersion, last.Version) == version {
			version = fmt.Sprintf(tplVersion, last.Version+1)
		}
	}
	filename = fmt.Sprintf(tplFileName, version, name, migrationType)
	fPath = filepath.Join(dir, filename)
	if tmpl = sqlMigrationTemplate; migrationType == migrationTypeGo {
		tmpl = goSQLMigrationTemplate
	}
	if migrationTemplate != nil {
		tmpl = migrationTemplate
	}
	if path, err = writeTemplateToFile(fPath, tmpl, version); err != nil {
		return
	}
	log.Printf(tplCreatedFile, path)

	return
}

// Create Создание нового пустого файла миграции.
func Create(db *sql.DB, dir, name, migrationType string) error {
	return CreateWithTemplate(db, dir, nil, name, migrationType)
}

// Запись шаблона миграции в файл.
func writeTemplateToFile(path string, t *template.Template, version string) (ret string, err error) {
	var f *os.File

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		err = Errors().FileAlreadyExists.Bind(path)
		return
	}
	if f, err = os.Create(path); err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	if err = t.Execute(f, version); err != nil {
		return
	}
	ret = f.Name()

	return
}
