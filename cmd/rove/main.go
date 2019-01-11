package main

import (
	"fmt"
	"os"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/database"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("rove", "Performs database migration tasks.")

	cDBPrefix  = app.Flag("envprefix", "Prefix for environment variables.").String()
	cDBAll     = app.Command("all", "Apply all changesets to the database.")
	cDBAllFile = cDBAll.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBUp      = app.Command("up", "Apply a specific number of changesets to the database.")
	cDBUpCount = cDBUp.Arg("count", "Number of changesets [int].").Required().Int()
	cDBUpFile  = cDBUp.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBReset     = app.Command("reset", "Run all rollbacks on the database.")
	cDBResetFile = cDBReset.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBDown      = app.Command("down", "Apply a specific number of rollbacks to the database.")
	cDBDownCount = cDBDown.Arg("count", "Number of rollbacks [int].").Required().Int()
	cDBDownFile  = cDBDown.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBStatus = app.Command("status", "Output the list of migrations already applied to the database.")
)

func main() {
	argList := os.Args[1:]
	arg := kingpin.MustParse(app.Parse(argList))

	// Create a new MySQL database object.
	conn, err := database.NewConnection(*cDBPrefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := database.NewMySQL(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch arg {
	case cDBAll.FullCommand():
		r := rove.NewFileMigration(m, *cDBAllFile)
		r.Verbose = true
		err = r.Migrate(0)
	case cDBUp.FullCommand():
		r := rove.NewFileMigration(m, *cDBUpFile)
		r.Verbose = true
		err = r.Migrate(*cDBUpCount)
	case cDBReset.FullCommand():
		r := rove.NewFileMigration(m, *cDBResetFile)
		r.Verbose = true
		err = r.Reset(0)
	case cDBDown.FullCommand():
		r := rove.NewFileMigration(m, *cDBDownFile)
		r.Verbose = true
		err = r.Reset(*cDBDownCount)
	case cDBStatus.FullCommand():
		r := rove.NewFileMigration(m, *cDBDownFile)
		r.Verbose = true
		_, err = r.Status()
	}

	// If there is an error, return with an error code of 1.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
