package main

import (
	"fmt"
	"os"

	"github.com/josephspurrier/rove"

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

	switch arg {
	case cDBAll.FullCommand():
		err := rove.Migrate(*cDBAllFile, *cDBPrefix, 0, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case cDBUp.FullCommand():
		err := rove.Migrate(*cDBUpFile, *cDBPrefix, *cDBUpCount, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case cDBReset.FullCommand():
		err := rove.Reset(*cDBResetFile, *cDBPrefix, 0, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case cDBDown.FullCommand():
		err := rove.Reset(*cDBDownFile, *cDBPrefix, *cDBDownCount, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
