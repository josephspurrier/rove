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

	assert.Contains(t, string(out), "Applied: 3) josephspurrier:3 (success.sql) 25f77f488ead6c37bc98ca17b3b52d81 [tag='']")
	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 0da13ad4c9f930f7cef760f2fea09854 [tag='']")
	assert.Contains(t, string(out), "Applied: 1) josephspurrier:1 (success.sql) a80573cbde97740723233be4ef760fe9 [tag='']")

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

	assert.Contains(t, string(out), "Applied: 1) josephspurrier:1 (success.sql) a80573cbde97740723233be4ef760fe9 [tag='']")
	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 0da13ad4c9f930f7cef760f2fea09854 [tag='']")

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

	assert.Contains(t, string(out), "Applied: 2) josephspurrier:2 (success.sql) 0da13ad4c9f930f7cef760f2fea09854 [tag='']")

	testutil.TeardownDatabase(unique)
}
