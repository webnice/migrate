package main

import "gopkg.in/alecthomas/kingpin.v2"

func args() (cmd string, args *Args) {
	args = new(Args)
	kingpin.CommandLine.Help = `Utility for applying database migrations`

	// Global flags
	kingpin.Flag(`dir`, `Directory with migration files. Default is '.'. Overrides the default value for a flag from an environment variable by name 'GOOSE_DIR'`).
		Envar(`GOOSE_DIR`).
		Default(".").
		StringVar(&args.Directory)
	kingpin.Flag(`drv`, `Driver of database. Support is [mysql, postgres, sqlite3, redshift, clickhouse, tidb]. Overrides the default value for a flag from an environment variable by name 'GOOSE_DRV'`).
		Envar("GOOSE_DRV").
		Default("mysql").
		StringVar(&args.Driver)
	kingpin.Flag(`dsn`, `Database source name (DSN). Overrides the default value for a flag from an environment variable by name 'GOOSE_DSN'`).
		Envar(`GOOSE_DSN`).
		Default(`root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true`).
		StringVar(&args.Dsn)

	// Commands with args
	args.Up = kingpin.Command(`up`, `Migrate the DB to the most recent version available`)
	args.UpTo = kingpin.Command(`up-to`, `Migrate the DB to a specific VERSION`)
	args.UpTo.Arg(`VERSION`, `Specific migration VERSION`).
		Default("").
		StringVar(&args.UpDownVersion)
	args.Down = kingpin.Command(`down`, `Roll back the version by 1`)
	args.DownTo = kingpin.Command(`down-to`, `Roll back to a specific VERSION`)
	args.DownTo.Arg(`VERSION`, `Specific migration VERSION`).
		Default("").
		StringVar(&args.UpDownVersion)
	args.Redo = kingpin.Command(`redo`, `Re-run the latest migration`)
	args.Status = kingpin.Command(`status`, `Dump the migration status for the current DB`)
	args.Version = kingpin.Command(`version`, `Print the current version of the database`)
	args.Create = kingpin.Command(`create`, `Creates new migration file with next version`)
	args.Create.Arg(`NAME`, `NAME of new migration file. Default name is 'new'`).
		Default("new").
		StringVar(&args.CreateName)
	args.Create.Arg(`TYPE`, `Type of new migration file [sql|go]. Default type is 'sql'`).
		Default("sql").
		StringVar(&args.CreateType)
	cmd = kingpin.Parse()

	return
}
