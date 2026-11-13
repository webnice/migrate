package goose

import (
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"log"
	"os"
	"strings"
)

// Проверка завершения строки точкой с запятой.
// Комментарий после двойного дефиса при проверке не учитывается.
func endsWithSemicolon(line string) (ret bool) {
	var (
		prev    string
		word    string
		scanner *bufio.Scanner
	)

	scanner = bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word = scanner.Text()
		if strings.HasPrefix(word, "--") {
			break
		}
		prev = word
	}
	ret = strings.HasSuffix(prev, ";")

	return
}

// Разбивает SQL-скрипт миграции на отдельные выражения.
//
// Базовый случай — разделение выражений точкой с запятой, естественно
// завершающей выражение. Более сложные случаи, например PL/pgSQL, содержат
// точки с запятой внутри выражения; для них предусмотрены явные аннотации
// StatementBegin и StatementEnd, отключающие разделение по точке с запятой.
//
// Вторым возвращаемым значением передаётся признак необходимости выполнения
// миграции внутри транзакции.
func getSQLStatements(r io.Reader, direction bool) (ret []string, tx bool, err error) {
	const (
		warnNoStatementEnd = "ВНИМАНИЕ: обнаружена аннотация '-- +goose StatementBegin' без парной '-- +goose StatementEnd'"
		tplUnfinished      = "ВНИМАНИЕ: незавершённое SQL-выражение: %s. Возможно пропущена точка с запятой"
	)
	var (
		buf               bytes.Buffer
		scanner           *bufio.Scanner
		line              string
		cmd               string
		upSections        int
		downSections      int
		statementEnded    bool
		ignoreSemicolons  bool
		directionIsActive bool
	)

	tx = true
	scanner = bufio.NewScanner(r)
	for scanner.Scan() {
		line = scanner.Text()
		// Обработка аннотаций goose.
		if strings.HasPrefix(line, sqlCmdPrefix) {
			cmd = strings.TrimSpace(line[len(sqlCmdPrefix):])
			switch cmd {
			case gooseUp:
				directionIsActive = direction
				upSections++
			case gooseDown:
				directionIsActive = !direction
				downSections++
			case gooseStatementBegin:
				if directionIsActive {
					ignoreSemicolons = true
				}
			case gooseStatementEnd:
				if directionIsActive {
					statementEnded = ignoreSemicolons
					ignoreSemicolons = false
				}
			case gooseNoTransaction:
				tx = false
			}
		}
		if !directionIsActive {
			continue
		}
		_, _ = buf.WriteString(line + "\n")
		// Завершение выражения, точка с запятой вне блока StatementBegin/End или явное окончание блока StatementEnd.
		if (!ignoreSemicolons && endsWithSemicolon(line)) || statementEnded {
			statementEnded = false
			ret = append(ret, buf.String())
			buf.Reset()
		}
	}
	if err = scanner.Err(); err != nil {
		err = Errors().ReadMigrationFile.Bind(err)
		return
	}
	// Диагностика вероятных ошибок скрипта миграции.
	if ignoreSemicolons {
		log.Println(warnNoStatementEnd)
	}
	if remaining := strings.TrimSpace(buf.String()); len(remaining) > 0 {
		log.Printf(tplUnfinished, remaining)
	}
	if upSections == 0 && downSections == 0 {
		err = Errors().NoAnnotations.Bind()
		return
	}

	return
}

// Выполнение миграции, записанной в виде SQL-скрипта.
//
// Секции скрипта размечаются специальным комментарием, начинающимся с
// "-- +goose", определяющим выполнение секции при прямой миграции (Up)
// или откате (Down). Все выражения после директивы направления объединяются
// в одну секцию до появления директивы противоположного направления.
func runSQLMigration(db *sql.DB, scriptFile string, v int64, direction bool) (err error) {
	var (
		f          *os.File
		statements []string
		useTx      bool
		tx         *sql.Tx
		query      string
		end        bool
	)

	if f, err = os.Open(scriptFile); err != nil {
		err = Errors().OpenMigrationFile.Bind(scriptFile, err)
		return
	}
	defer func() { _ = f.Close() }()
	if statements, useTx, err = getSQLStatements(f, direction); err != nil {
		return
	}
	if useTx {
		if tx, err = db.Begin(); err != nil {
			err = Errors().BeginTransaction.Bind(err)
			return
		}
		for _, query = range statements {
			if end {
				break
			}
			if _, err = tx.Exec(query); err != nil {
				err = Errors().ExecStatement.Bind(err)
				end = true
			}
		}
		if !end {
			// Запись версии миграции в лог версий.
			if _, err = tx.Exec(GetDialect().insertVersionSQL(), v, direction); err != nil {
				err = Errors().InsertVersion.Bind(err)
				end = true
			}
		}
		if end {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}
	// Выполнение без транзакции.
	for _, query = range statements {
		if end {
			break
		}
		if _, err = db.Exec(query); err != nil {
			err = Errors().ExecStatement.Bind(err)
			end = true
		}
	}
	if end {
		return
	}
	if _, err = db.Exec(GetDialect().insertVersionSQL(), v, direction); err != nil {
		err = Errors().InsertVersion.Bind(err)
	}

	return
}
