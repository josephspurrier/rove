package jsonfile_test

import (
	"os"
	"testing"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/jsonfile"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/josephspurrier/rove/pkg/adapter/mysql/testutil"

	"github.com/stretchr/testify/assert"
)

func TestFileMigration(t *testing.T) {
	_, unique := testutil.SetupDatabase()
	_ = os.Remove("test.json")

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	js, err := jsonfile.New("test.json", m)
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(js, "../testdata/success.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:3", s)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:2", s)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:1", s)

	_ = os.Remove("test.json")
	testutil.TeardownDatabase(unique)
}

func TestMigrationFailDuplicate(t *testing.T) {
	_, unique := testutil.SetupDatabase()
	_ = os.Remove("test.json")

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	js, err := jsonfile.New("test.json", m)
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(js, "../testdata/fail-duplicate.sql")
	r.Verbose = true

	err = r.Migrate(0)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "checksum does not match")

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:2", s)

	_ = os.Remove("test.json")
	testutil.TeardownDatabase(unique)
}

func TestInclude(t *testing.T) {
	_, unique := testutil.SetupDatabase()
	_ = os.Remove("test.json")

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	js, err := jsonfile.New("test.json", m)
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(js, "../testdata/parent.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:3", s)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:2", s)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:1", s)

	_ = os.Remove("test.json")
	testutil.TeardownDatabase(unique)
}

func TestChangesetMigration(t *testing.T) {
	_, unique := testutil.SetupDatabase()
	_ = os.Remove("test.json")

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	js, err := jsonfile.New("test.json", m)
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewChangesetMigration(js, sSuccess)
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:3", s)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:2", s)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "josephspurrier:1", s)

	_ = os.Remove("test.json")
	testutil.TeardownDatabase(unique)
}

var sSuccess = `
--changeset josephspurrier:1
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user_status;

--changeset josephspurrier:2
INSERT INTO user_status (id, status, created_at, updated_at, deleted) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);
--rollback TRUNCATE TABLE user_status;

--changeset josephspurrier:3
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
CREATE TABLE user (
    id VARCHAR(36) NOT NULL,
    
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    
    status_id TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT 0,
    
    UNIQUE KEY (email),
    CONSTRAINT f_user_status FOREIGN KEY (status_id) REFERENCES user_status (id) ON DELETE CASCADE ON UPDATE CASCADE,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user;`
