package rove_test

import (
	"testing"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/josephspurrier/rove/pkg/adapter/mysql/testutil"

	"github.com/stretchr/testify/assert"
)

func TestFileLoadFailure(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/success.sql")
	r.Verbose = true

	for v := range []error{
		r.Migrate(0),
		r.Reset(0),
		r.Rollback("none"),
		func() error {
			_, err := r.Status()
			return err
		}(),
		r.Tag("none"),
	} {
		assert.NotNil(t, v)
	}
}

func TestFileLoadFailure2(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "not-exist")
	r.Verbose = true

	for v := range []error{
		r.Migrate(0),
		r.Reset(0),
		r.Rollback("none"),
		func() error {
			_, err := r.Status()
			return err
		}(),
		r.Tag("none"),
	} {
		assert.NotNil(t, v)
	}
}

func TestStringLoadFailure(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewChangesetMigration(m, sFail)
	r.Verbose = true

	for v := range []error{
		r.Migrate(0),
		r.Reset(0),
		r.Rollback("none"),
		func() error {
			_, err := r.Status()
			return err
		}(),
		r.Tag("none"),
	} {
		assert.NotNil(t, v)
	}
}

func TestLoadChangesetFailure(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/fail-duplicate.sql")
	r.Verbose = true

	for v := range []error{
		r.Migrate(0),
		r.Reset(0),
		r.Rollback("none"),
		func() error {
			_, err := r.Status()
			return err
		}(),
		r.Tag("none"),
	} {
		assert.NotNil(t, v)
	}
}

func TestLoadChanglogFailure(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/success.sql")
	r.Verbose = true

	m.InitializeQuery = "INSERT bad query;"

	for v := range []error{
		r.Migrate(0),
		r.Reset(0),
		r.Rollback("none"),
		func() error {
			_, err := r.Status()
			return err
		}(),
		r.Tag("none"),
	} {
		assert.NotNil(t, v)
	}
}

func TestFileMigration(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/success.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "3", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Nil(t, s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Show status of the migrations.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	testutil.TeardownDatabase(unique)
}

func TestMigrationFailDuplicate(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/fail-duplicate.sql")
	r.Verbose = true

	err = r.Migrate(0)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "checksum does not match")

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	testutil.TeardownDatabase(unique)
}

func TestInclude(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/parent.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "3", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Nil(t, s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	testutil.TeardownDatabase(unique)
}

func TestFileBadQuery(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/badquery.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.NotNil(t, err)
}

func TestFileBadSpace(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/badspace.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.NotNil(t, err)
}

func TestFileMissingHeader(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/missingheader.sql")
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.NotNil(t, err)
}

func TestChangesetMigration(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewChangesetMigration(m, sSuccess)
	r.Verbose = true

	// Run migration.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "3", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Run migration again.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Remove all migrations.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Nil(t, s)

	// Remove all migrations again.
	err = r.Reset(0)
	assert.Nil(t, err)

	// Run 2 migrations.
	err = r.Migrate(2)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	// Remove 1 migration.
	err = r.Reset(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	testutil.TeardownDatabase(unique)
}

func TestChangesetTag(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewChangesetMigration(m, sSuccess)
	r.Verbose = true

	// Run migration.
	err = r.Migrate(1)
	assert.Nil(t, err)

	// Tag the migration.
	err = r.Tag("jas1")
	assert.Nil(t, err)

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)
	assert.Equal(t, "jas1", s.Tag)

	// Run migration again.
	err = r.Migrate(1)
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)
	assert.Equal(t, "", s.Tag)

	// Rollback to the tag.
	err = r.Rollback("jas1")
	assert.Nil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)
	assert.Equal(t, "jas1", s.Tag)

	// Attempt rollback again.
	err = r.Rollback("jas1")
	assert.NotNil(t, err)

	// Get the status.
	s, err = r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "1", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)
	assert.Equal(t, "jas1", s.Tag)

	// Run migration again.
	err = r.Migrate(1)
	assert.Nil(t, err)

	// Attempt to tag with the same tag.
	err = r.Tag("jas1")
	assert.NotNil(t, err)

	// Attempt to tag with an empty string.
	err = r.Tag("")
	assert.NotNil(t, err)

	// Attempt rollback to a tag that doesn't exist.
	err = r.Rollback("not-exist")
	assert.NotNil(t, err)

	// Attempt rollback to an empty tag.
	err = r.Rollback("")
	assert.NotNil(t, err)

	testutil.TeardownDatabase(unique)
}

var sSuccess = `
--changeset josephspurrier:1
--description Create a user_status table.
-- Standard comment.
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

var sFail = `
--description Create a user_status table.
-- Standard comment.
CREATE TABLE user_status (
--rollback DROP TABLE user_status;`
