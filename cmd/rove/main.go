package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	checksumError  = "error"
	checksumIgnore = "ignore"
	checksumUpdate = "update"
)

var (
	app       = kingpin.New("rove", "Performs database migration tasks.")
	cChecksum = app.Flag("checksum-mode", "Set how to handle checksums that don't match "+
		"[error (default),ignore,update].").Enum(checksumError, checksumIgnore, checksumUpdate)

	cDBHost      = app.Flag("hostname", "Database hostname or IP [string].").String()
	cDBPort      = app.Flag("port", "Database port [int].").Int()
	cDBUsername  = app.Flag("username", "Database username [string].").String()
	cDBPassword  = app.Flag("password", "Database password [string].").String()
	cDBName      = app.Flag("name", "Database name [string].").String()
	cDBParameter = app.Flag("parameter", "Database parameters [string].").String()

	cDBPrefix  = app.Flag("envprefix", "Prefix for environment variables.").String()
	cDBAll     = app.Command("all", "Apply all changesets to the database.")
	cDBAllFile = cDBAll.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBUp      = app.Command("up", "Apply a specific number of changesets to the database.")
	cDBUpCount = cDBUp.Arg("count", "Number of changesets [int].").Required().Int()
	cDBUpFile  = cDBUp.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBReset     = app.Command("reset", "Apply all rollbacks to the database.")
	cDBResetFile = cDBReset.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBDown      = app.Command("down", "Apply a specific number of rollbacks to the database.")
	cDBDownCount = cDBDown.Arg("count", "Number of rollbacks [int].").Required().Int()
	cDBDownFile  = cDBDown.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBTag     = app.Command("tag", "Apply a tag to the latest changeset in the database.")
	cDBTagName = cDBTag.Arg("name", "Name of the tag [string].").Required().String()
	cDBTagFile = cDBTag.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBRollback     = app.Command("rollback", "Run all rollbacks until the specified tag on the database.")
	cDBRollbackName = cDBRollback.Arg("name", "Name of the tag [string].").Required().String()
	cDBRollbackFile = cDBRollback.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBConvert     = app.Command("convert", "Convert a Liquibase changelog table to a Rove changelog table.")
	cDBConvertFile = cDBConvert.Arg("file", "Filename of the migration file [string].").Required().String()

	cDBStatus = app.Command("status", "Output the list of changesets already applied to the database.")
)

func main() {
	argList := os.Args[1:]
	arg := kingpin.MustParse(app.Parse(argList))

	// Set the ChecksumMode.
	csMode := rove.ChecksumThrowError
	if *cChecksum == checksumIgnore {
		csMode = rove.ChecksumIgnore
	} else if *cChecksum == checksumUpdate {
		csMode = rove.ChecksumUpdate
	}

	// Create the MySQL connection information from environment variables.
	conn, err := mysql.NewConnection(*cDBPrefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Overwrite the database variables if there are parameters set.
	if len(*cDBHost) > 0 {
		conn.Hostname = *cDBHost
	}
	if *cDBPort > 0 {
		conn.Port = *cDBPort
	}
	if len(*cDBUsername) > 0 {
		conn.Username = *cDBUsername
	}
	if len(*cDBPassword) > 0 {
		conn.Password = *cDBPassword
	}
	if len(*cDBName) > 0 {
		conn.Name = *cDBName
	}
	if len(*cDBParameter) > 0 {
		conn.Parameter = *cDBParameter
	}

	// Add parseTime parameter if it's not included to parse times properly
	// in MySQL.
	if !strings.Contains(conn.Parameter, "parseTime") {
		if conn.Parameter == "" {
			conn.Parameter = "parseTime=true"
		} else {
			conn.Parameter += "&parseTime=true"
		}
	}

	// Create a new MySQL database object.
	db, err := mysql.New(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch arg {
	case cDBAll.FullCommand():
		r := rove.NewFileMigration(db, *cDBAllFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Migrate(0)
	case cDBUp.FullCommand():
		r := rove.NewFileMigration(db, *cDBUpFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Migrate(*cDBUpCount)
	case cDBReset.FullCommand():
		r := rove.NewFileMigration(db, *cDBResetFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Reset(0)
	case cDBDown.FullCommand():
		r := rove.NewFileMigration(db, *cDBDownFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Reset(*cDBDownCount)
	case cDBTag.FullCommand():
		r := rove.NewFileMigration(db, *cDBTagFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Tag(*cDBTagName)
	case cDBRollback.FullCommand():
		r := rove.NewFileMigration(db, *cDBRollbackFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Rollback(*cDBRollbackName)
	case cDBConvert.FullCommand():
		r := rove.NewFileMigration(db, *cDBConvertFile)
		r.Verbose = true
		r.Checksum = csMode
		err = r.Convert(db.DB)
	case cDBStatus.FullCommand():
		r := rove.NewFileMigration(db, "")
		r.Verbose = true
		r.Checksum = csMode
		_, err = r.Status()
	}

	// If there is an error, return with an error code of 1.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
