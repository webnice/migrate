package main

import (
//	"os/exec"
//	"strings"
//	"testing"
)

//func TestDefaultBinary(t *testing.T) {
//	commands := []string{
//		"go build -i -o ../bin/gsmigrate github.com/webnice/migrate/gsmigrate",
//		//"./gsmigrate -dir=examples/sql-migrations sqlite3 sql.db up",
//		//"./gsmigrate -dir=examples/sql-migrations sqlite3 sql.db version",
//		//"./gsmigrate -dir=examples/sql-migrations sqlite3 sql.db down",
//		"../bin/gsmigrate --drv=sqlite3 --dsn=../sql.db status",
//	}
//
//	for _, cmd := range commands {
//		args := strings.Split(cmd, " ")
//		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
//		if err != nil {
//			t.Fatalf("%s:\n%v\n\n%s", err, cmd, out)
//		}
//	}
//}

//func TestCustomBinary(t *testing.T) {
//	commands := []string{
//		"go build -i -o custom-goose ./examples/go-migrations",
//		"./custom-goose -dir=examples/go-migrations sqlite3 go.db up",
//		"./custom-goose -dir=examples/go-migrations sqlite3 go.db version",
//		"./custom-goose -dir=examples/go-migrations sqlite3 go.db down",
//		"./custom-goose -dir=examples/go-migrations sqlite3 go.db status",
//	}
//
//	for _, cmd := range commands {
//		args := strings.Split(cmd, " ")
//		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
//		if err != nil {
//			t.Fatalf("%s:\n%v\n\n%s", err, cmd, out)
//		}
//	}
//}
