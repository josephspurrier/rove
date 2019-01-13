package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/josephspurrier/rove/pkg/adapter/mysql/testutil"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMigrationAll(t *testing.T) {
	_, unique := migrateAll(t)
	testutil.TeardownDatabase(unique)
}

func migrateAll(t *testing.T) (*sqlx.DB, string) {
	db, unique := testutil.SetupDatabase()

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
	testutil.SetEnv(unique)
	main()
	testutil.UnsetEnv(unique)

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "josephspurrier:1")
	assert.Contains(t, string(out), "josephspurrier:2")
	assert.Contains(t, string(out), "josephspurrier:3")

	return db, unique
}

func TestMigrationReset(t *testing.T) {
	_, unique := migrateAll(t)

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
	testutil.SetEnv(unique)
	main()
	testutil.UnsetEnv(unique)

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Applied: 3) josephspurrier:3 (success.sql) 57cc0b1c45cb72032bcaed07483d243d [tag='']")
	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 3f81b08751b27ff6680b287e08ea112a [tag='']")
	assert.Contains(t, string(out), "Applied: 1) josephspurrier:1 (success.sql) f0685b81bd072358e0ea0cd0750b544a [tag='']")

	testutil.TeardownDatabase(unique)
}

func TestMigrationUp(t *testing.T) {
	_, unique := migrateUp(t)
	testutil.TeardownDatabase(unique)
}

func migrateUp(t *testing.T) (*sqlx.DB, string) {
	db, unique := testutil.SetupDatabase()

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
	testutil.SetEnv(unique)
	main()
	testutil.UnsetEnv(unique)

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Applied: 1) josephspurrier:1 (success.sql) f0685b81bd072358e0ea0cd0750b544a [tag='']")
	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 3f81b08751b27ff6680b287e08ea112a [tag='']")

	return db, unique
}

func TestMigrationDown(t *testing.T) {
	_, unique := migrateUp(t)

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
	testutil.SetEnv(unique)
	main()
	testutil.UnsetEnv(unique)

	// Get the output.
	w.Close()
	out, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	os.Stdout = backupd

	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 3f81b08751b27ff6680b287e08ea112a [tag='']")

	testutil.TeardownDatabase(unique)
}
