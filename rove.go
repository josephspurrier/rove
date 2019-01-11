package rove

import (
	"errors"
)

const (
	appVersion       = "1.0"
	elementChangeset = "--changeset "
	elementRollback  = "--rollback "
	elementInclude   = "--include "
	elementMemory    = "memory"
)

var (
	// ErrInvalidFormat is when a changeset is not found.
	ErrInvalidFormat = errors.New("invalid changeset format")
)

// Rove contains the database migration information.
type Rove struct {
	// Verbose is whether information is written to the screen or not.
	Verbose bool

	// file is the full path to the migration file.
	file string
	// changeset is text with changesets.
	changeset string
	// db is a migration.
	db Migration
}

// NewFileMigration returns a file migration object.
func NewFileMigration(db Migration, filename string) *Rove {
	return &Rove{
		db:   db,
		file: filename,
	}
}

// NewChangesetMigration returns a changeset migration object.
func NewChangesetMigration(db Migration, changeset string) *Rove {
	return &Rove{
		db:        db,
		changeset: changeset,
	}
}
