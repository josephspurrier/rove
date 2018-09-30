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
	// m is a migration.
	m Migration
}

// DBChangeset contains a single database record change.
type DBChangeset struct {
	ID            string
	Author        string
	Filename      string
	OrderExecuted int
}

// New returns an instance of Rove.
func New(m Migration) *Rove {
	return &Rove{
		m: m,
	}
}
