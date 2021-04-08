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

// Create writes a new blank migration file.
func CreateWithTemplate(db *sql.DB, dir string, migrationTemplate *template.Template, name, migrationType string) (err error) {
	var migrations Migrations
	var last *Migration
	var version, filename, fpath, path string
	var tmpl *template.Template

	migrations, err = CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return
	}
	version = time.Now().In(time.UTC).Format("20060102150405")
	if last, err = migrations.Last(); err == nil {
		if fmt.Sprintf("%014d", last.Version) == version {
			version = fmt.Sprintf("%014d", last.Version+1)
		}
	}
	filename = fmt.Sprintf("%v_%v.%v", version, name, migrationType)
	fpath = filepath.Join(dir, filename)
	tmpl = sqlMigrationTemplate
	if migrationType == "go" {
		tmpl = goSQLMigrationTemplate
	}
	if migrationTemplate != nil {
		tmpl = migrationTemplate
	}
	path, err = writeTemplateToFile(fpath, tmpl, version)
	if err != nil {
		return
	}
	log.Printf("Created new file: %s\n", path)

	return
}

// Create writes a new blank migration file.
func Create(db *sql.DB, dir, name, migrationType string) error {
	return CreateWithTemplate(db, dir, nil, name, migrationType)
}

func writeTemplateToFile(path string, t *template.Template, version string) (string, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to create file: %v already exists", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = t.Execute(f, version)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

var sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(
	`-- +goose Up
-- SQL in this section is executed when the migration is applied.



-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

`))

var goSQLMigrationTemplate = template.Must(template.New("goose.go-migration").Parse(
	`package main

import (
	"database/sql"

	"gopkg.in/webnice/migrate.v1/goose"
)

func init() {
	goose.AddMigration(Up{{.}}, Down{{.}})
}

// Up{{.}} Migration applied
func Up{{.}}(tx *sql.Tx) (err error) {
	// This code is executed when the migration is applied.
	return
}

// Down{{.}} Migration rolled back
func Down{{.}}(tx *sql.Tx) (err error) {
	// This code is executed when the migration is rolled back.
	return
}
`))
