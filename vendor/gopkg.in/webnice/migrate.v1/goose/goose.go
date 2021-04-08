package goose

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

var (
	duplicateCheckOnce sync.Once
	minVersion         = int64(0)
	maxVersion         = int64((1 << 63) - 1)
)

// Run runs a goose command.
func Run(command string, db *sql.DB, dir string, args ...string) (err error) {
	var version int64

	switch command {
	case "up":
		err = Up(db, dir)
	case "up-by-one":
		err = UpByOne(db, dir)
	case "up-to":
		if len(args) == 0 {
			err = fmt.Errorf("up-to must be of form: goose [OPTIONS] DRIVER DBSTRING up-to VERSION")
			return
		}
		version, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("version must be a number (got '%s')", args[0])
			return
		}
		err = UpTo(db, dir, version)
	case "create":
		if len(args) == 0 {
			err = fmt.Errorf("create must be of form: goose [OPTIONS] DRIVER DBSTRING create NAME [go|sql]")
			return
		}
		migrationType := "go"
		if len(args) == 2 {
			migrationType = args[1]
		}
		err = Create(db, dir, args[0], migrationType)
	case "down":
		err = Down(db, dir)
	case "down-to":
		if len(args) == 0 {
			err = fmt.Errorf("down-to must be of form: goose [OPTIONS] DRIVER DBSTRING down-to VERSION")
			return
		}
		version, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("version must be a number (got '%s')", args[0])
			return
		}
		err = DownTo(db, dir, version)
	case "redo":
		err = Redo(db, dir)
	case "reset":
		err = Reset(db, dir)
	case "status":
		err = Status(db, dir)
	case "version":
		err = Version(db, dir)
	default:
		err = fmt.Errorf("%q: no such command", command)
	}

	return
}
