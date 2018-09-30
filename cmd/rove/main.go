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

	cDB        = app.Command("migrate", "Perform actions on the database.")
	cDBPrefix  = cDB.Flag("envprefix", "Prefix for environment variables.").String()
	cDBAll     = cDB.Command("all", "Apply all changesets to the database.")
	cDBAllFile = cDBAll.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBUp      = cDB.Command("up", "Apply a specific number of changesets to the database.")
	cDBUpCount = cDBUp.Arg("count", "Number of changesets [int].").Required().Int()
	cDBUpFile  = cDBUp.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBReset     = cDB.Command("reset", "Run all rollbacks on the database.")
	cDBResetFile = cDBReset.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBDown      = cDB.Command("down", "Apply a specific number of rollbacks to the database.")
	cDBDownCount = cDBDown.Arg("count", "Number of rollbacks [int].").Required().Int()
	cDBDownFile  = cDBDown.Arg("file", "Filename of the migration file [string].").Required().String()
)

func main() {
	argList := os.Args[1:]
	arg := kingpin.MustParse(app.Parse(argList))

	// Create a new MySQL database object.
	m, err := database.NewMySQL(*cDBPrefix)
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
	}

	// If there is an error, return with an error code of 1.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
