package rove_test

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/josephspurrier/rove/pkg/adapter/mysql/testutil"

	"github.com/stretchr/testify/assert"
)

func TestFileLoadFailure(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	for _, f := range []func() *rove.Rove{
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/success.sql")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "not-exist")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Read the file into a string.
			b, err := ioutil.ReadFile("testdata/missingheader.sql")
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewChangesetMigration(m, string(b))
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/fail-duplicate.sql")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/success.sql")
			r.Verbose = true
			m.InitializeQuery = "INSERT bad query;"
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/badquery.sql")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/badspace.sql")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/missingheader.sql")
			r.Verbose = true
			return r
		},
		func() *rove.Rove {
			con := testutil.Connection(unique)
			con.Username = "wrong"

			// Create a new MySQL database object.
			m, err := mysql.New(con)
			assert.NotNil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/missingheader.sql")
			r.Verbose = true
			return r
		},
	} {
		rr := f()
		for v := range []error{
			rr.Migrate(0),
			rr.Reset(0),
			rr.Rollback("none"),
			func() error {
				_, err := rr.Status()
				return err
			}(),
			rr.Tag("none"),
		} {
			assert.NotNil(t, v)
		}
	}

	testutil.TeardownDatabase(unique)
}

func TestFileMigration(t *testing.T) {
	for _, f := range []func() (*rove.Rove, string){
		func() (*rove.Rove, string) {
			_, unique := testutil.SetupDatabase()

			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewFileMigration(m, "testdata/success.sql")
			r.Verbose = true
			return r, unique
		},
		func() (*rove.Rove, string) {
			_, unique := testutil.SetupDatabase()

			// Create a new MySQL database object.
			m, err := mysql.New(testutil.Connection(unique))
			assert.Nil(t, err)

			// Read the file into a string.
			b, err := ioutil.ReadFile("testdata/success.sql")
			assert.Nil(t, err)

			// Set up rove.
			r := rove.NewChangesetMigration(m, string(b))
			r.Verbose = true
			return r, unique
		},
	} {
		r, unique := f()

		// Run migration.
		err := r.Migrate(0)
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
	assert.Contains(t, err.Error(), "duplicate entry found")

	// Get the status.
	s, err := r.Status()
	assert.Nil(t, err)
	assert.Equal(t, "2", s.ID)
	assert.Equal(t, "josephspurrier", s.Author)

	testutil.TeardownDatabase(unique)
}

func TestMigrationFailChecksum(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/success.sql")
	r.Verbose = true

	// Perform all migrations.
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Change a checksum.
	err = m.Update("1", "josephspurrier", "success.sql", time.Now(), 1,
		"bad", "description", "version")
	assert.Nil(t, err)

	// Migrate and throw error.
	r.Checksum = rove.ChecksumThrowError
	err = r.Migrate(0)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "checksum does not match")

	// Migrate and ignore error.
	r.Checksum = rove.ChecksumIgnore
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Migrate and throw error.
	r.Checksum = rove.ChecksumThrowError
	err = r.Migrate(0)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "checksum does not match")

	// Migrate and update checksum.
	r.Checksum = rove.ChecksumUpdate
	err = r.Migrate(0)
	assert.Nil(t, err)

	// Migrate and no longer throw an error.
	r.Checksum = rove.ChecksumThrowError
	err = r.Migrate(0)
	assert.Nil(t, err)

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
func TestChangesetTag(t *testing.T) {
	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Read the file into a string.
	b, err := ioutil.ReadFile("testdata/success.sql")
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewChangesetMigration(m, string(b))
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

func TestTransactions(t *testing.T) {

	_, unique := testutil.SetupDatabase()

	// Create a new MySQL database object.
	m, err := mysql.New(testutil.Connection(unique))
	assert.Nil(t, err)

	// Set up rove.
	r := rove.NewFileMigration(m, "testdata/success.sql")
	r.Verbose = true

	err = r.Migrate(1)
	assert.Nil(t, err)

	// Clear the transaction func.
	m.TransactionFunc = nil

	// Fail on BeginTx().
	err = r.Migrate(1)
	assert.NotNil(t, err)
	err = r.Reset(1)
	assert.NotNil(t, err)

	// Add the transaction mock.
	mock := new(TxMock)
	m.TransactionFunc = func(tx *sql.Tx) rove.Transaction {
		return mock
	}

	// Fail on commit.
	mock.CommitError = errors.New("error")
	err = r.Migrate(1)
	assert.NotNil(t, err)
	err = r.Reset(1)
	assert.NotNil(t, err)

	// Fail on rollback.
	mock.RollbackError = errors.New("error")
	err = r.Migrate(1)
	assert.NotNil(t, err)
	err = r.Reset(1)
	assert.NotNil(t, err)
	mock.Reset()

	// Fail on exec.
	mock.ExecError = errors.New("error")
	err = r.Migrate(1)
	assert.NotNil(t, err)
	err = r.Reset(1)
	assert.NotNil(t, err)
	mock.Reset()

	testutil.TeardownDatabase(unique)
}

// TxMock is a database transaction mock.
type TxMock struct {
	CommitError   error
	RollbackError error
	ExecError     error
}

// Reset will reset all the errors to nil.
func (t *TxMock) Reset() {
	t.CommitError = nil
	t.RollbackError = nil
	t.ExecError = nil
}

// Commit will commit changes to the database or return an error.
func (t *TxMock) Commit() error {
	return t.CommitError
}

// Rollback will rollback changes to the database or return an error.
func (t *TxMock) Rollback() error {
	return t.RollbackError
}

// Exec will run a query on the database.
func (t *TxMock) Exec(query string) error {
	return t.ExecError
}
