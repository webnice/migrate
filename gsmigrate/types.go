package main

import "gopkg.in/alecthomas/kingpin.v2"

type Args struct {
	Up            *kingpin.CmdClause // Migrate the DB to the most recent version available
	UpTo          *kingpin.CmdClause // Migrate the DB to a specific VERSION
	Down          *kingpin.CmdClause // Roll back the version by 1
	DownTo        *kingpin.CmdClause // Roll back to a specific VERSION
	UpDownVersion string             // Specific migration VERSION
	Redo          *kingpin.CmdClause // Re-run the latest migration
	Status        *kingpin.CmdClause // Dump the migration status for the current DB
	Version       *kingpin.CmdClause // Print the current version of the database
	Create        *kingpin.CmdClause // Creates new migration file with next version
	CreateName    string             // NAME of new migration file
	CreateType    string             // Type of new migration file [sql|go]
	Directory     string             // Directory with migration files
	Driver        string             // Driver of database
	Dsn           string             // Database source name (DSN)
}
