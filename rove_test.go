package rove_test

import (
	"testing"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/database"
	"github.com/josephspurrier/rove/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	db, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m := new(database.MySQL)
	m.DB = db

	// Set up rove.
	r := rove.New(m)
	r.MigrationFile = "testdata/success.sql"

	// Run migration.
	err := r.Migrate(0)
	assert.Nil(t, err)

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 3, rows)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 0, rows)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 2, rows)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 1, rows)

	testutil.TeardownDatabase(unique)
}

func TestMigrationFailDuplicate(t *testing.T) {
	db, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m := new(database.MySQL)
	m.DB = db

	// Set up rove.
	r := rove.New(m)
	r.MigrationFile = "testdata/fail-duplicate.sql"

	err := r.Migrate(0)
	assert.Contains(t, err.Error(), "checksum does not match")

	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 2, rows)

	testutil.TeardownDatabase(unique)
}

func TestInclude(t *testing.T) {
	db, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m := new(database.MySQL)
	m.DB = db

	// Set up rove.
	r := rove.New(m)
	r.MigrationFile = "testdata/parent.sql"

	// Run migration.
	err := r.Migrate(0)
	assert.Nil(t, err)

	// Count the records.
	rows := 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 3, rows)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 0, rows)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 2, rows)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	rows = 0
	err = db.Get(&rows, `SELECT count(*) from databasechangelog`)
	assert.Nil(t, err)
	assert.Equal(t, 1, rows)

	testutil.TeardownDatabase(unique)
}
