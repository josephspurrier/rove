package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/josephspurrier/rove/pkg/adapter/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMigrationAll(t *testing.T) {
	_, unique := migrateAll(t)
	mysql.TeardownDatabase(unique)
}

func migrateAll(t *testing.T) (*sqlx.DB, string) {
	db, unique := mysql.SetupDatabase()

	// Set the arguments.
	os.Args = []string{
		"rove",
		"all",
		"testdata/success.sql",
		"--envprefix",
		unique,
	}

	// Redirect stdout.
	backupd := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the application.
	main()

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Changeset applied")

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 3, rows)

	return db, unique
}

func TestMigrationReset(t *testing.T) {
	db, unique := migrateAll(t)

	// Set the arguments.
	os.Args = []string{
		"rove",
		"reset",
		"testdata/success.sql",
		"--envprefix",
		unique,
	}

	// Redirect stdout.
	backupd := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the application.
	main()

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Rollback applied")

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 0, rows)

	mysql.TeardownDatabase(unique)
}

func TestMigrationUp(t *testing.T) {
	_, unique := migrateUp(t)
	mysql.TeardownDatabase(unique)
}

func migrateUp(t *testing.T) (*sqlx.DB, string) {
	db, unique := mysql.SetupDatabase()

	// Set the arguments.
	os.Args = []string{
		"rove",
		"up",
		"2",
		"testdata/success.sql",
		"--envprefix",
		unique,
	}

	// Redirect stdout.
	backupd := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the application.
	main()

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Changeset applied")

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 2, rows)

	return db, unique
}

func TestMigrationDown(t *testing.T) {
	db, unique := migrateUp(t)

	// Set the arguments.
	os.Args = []string{
		"rove",
		"down",
		"1",
		"testdata/success.sql",
		"--envprefix",
		unique,
	}

	// Redirect stdout.
	backupd := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the application.
	main()

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Rollback applied")

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 1, rows)

	mysql.TeardownDatabase(unique)
}
