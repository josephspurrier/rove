package rove

import (
	"database/sql"
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
	DatabaseType string
	// m is a migration.
	m Migration
}

// Migration represents functions to run on a database.
type Migration interface {
	// CreateChangelogTable returns an error if there is a problem with the
	// query. If it already exists, it should return nil.
	CreateChangelogTable() error
	// ChangesetApplied should return an error if there is a problem with the
	// query. If there are no rows, that error can be returned and it will
	// be ignored.
	ChangesetApplied(id string, author string, filename string) (checksum string, err error)
	// BeginTx starts a transaction.
	BeginTx() (*sql.Tx, error)
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

// New returns an instance of Rove.
func New(m Migration) *Rove {
	return &Rove{
		m: m,
	}
}
