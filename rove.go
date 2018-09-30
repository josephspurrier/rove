package rove

import (
	"errors"
)

const (
	appVersion       = "1.0"
	elementChangeset = "--changeset "
	elementRollback  = "--rollback "
	elementInclude   = "--include "
)

var (
	// ErrInvalidHeader is when the changeset header is invalid.
	ErrInvalidHeader = errors.New("invalid changeset header")
	// ErrInvalidFormat is when a changeset is not found.
	ErrInvalidFormat = errors.New("invalid changeset format")
)

// Rove contains the database migration information.
type Rove struct {
	// Verbose is whether information is written to the screen or not.
	Verbose bool
	// MigrationFile is the full path to the migration file.
	MigrationFile string
	// DatabaseType can be: mysql.
	//DatabaseType string
	// m is a migration.
	m Migration
}

// DBChangeset represents the database table records.
type DBChangeset struct {
	ID            string `db:"id"`
	Author        string `db:"author"`
	Filename      string `db:"filename"`
	OrderExecuted int    `db:"orderexecuted"`
}

// Migration represents a database specific migration.
type Migration interface {
	// CreateChangelogTable returns an error if there is a problem with the
	// query. If it already exists, it should return nil.
	CreateChangelogTable() error
	// ChangesetApplied should return an error if there is a problem with the
	// query. If there are no rows, that error can be returned and it will
	// be ignored.
	ChangesetApplied(id string, author string, filename string) (checksum string, err error)
	// BeginTx starts a transaction.
	BeginTx() (Transaction, error)
	// Count returns the number of changesets in the database and returns an
	// error if there is a problem with the query.
	Count() (count int, err error)
	// Insert will insert a new record into the database.
	Insert(id, author, filename string, count int, checksum, description, version string) error
	// Returns a list of the changesets or it returns an error if there is an
	// problem running the query.
	Changesets() ([]DBChangeset, error)
	// Delete the changeset from the database or return an error if there is a
	// problem running the query.
	Delete(id, author, filename string) error
}

// Transaction represents a database transaction.
type Transaction interface {
	// Commit will attempt to commit the changes to the database or return
	// an error.
	Commit() error
	// Rollback rollback changes to the database after a filed commit.
	Rollback() error
	// Exec runs a query on the database.
	Exec(query string, args ...interface{}) error
}

// New returns an instance of Rove.
func New(m Migration) *Rove {
	return &Rove{
		m: m,
	}
}
