package rove

import (
	"time"

	"github.com/josephspurrier/rove/pkg/changeset"
)

// Changelog represents a list of operations for a changelog.
type Changelog interface {
	// Initialize should perform any work to set up the changelog or return an
	// error.
	Initialize() error
	// BeginTx should start a transaction on the changelog.
	BeginTx() (Transaction, error)
	// ChangesetApplied should return the checksum from the changelog of a
	// matching changeset, an error, or a blank string if the changeset doesn't
	// exist in the changelog.
	ChangesetApplied(id string, author string, filename string) (record *changeset.Record, err error)
	// Changesets should return a list of the changesets or return an error.
	Changesets(reverse bool) ([]changeset.Record, error)
	// Count should return the number of changesets in the changelog or return
	// an error.
	Count() (count int, err error)
	// Insert should add a new changeset to the changelog or return an error.
	Insert(id, author, filename string, dateexecuted time.Time, count int,
		checksum, description, version string) error
	// Update should update a changeset in the changelog or return an error.
	Update(id, author, filename string, dateexecuted time.Time, count int,
		checksum, description, version string) error
	// Delete should remove the changeset from the changelog or return an error.
	Delete(id, author, filename string) error
	// Tag should add a tag to the latest changeset in the database or return
	// an error.
	Tag(id, author, filename, tag string) error
	// Rollback should return the number of changests to remove to revert to the
	// specified tag or return an error.
	Rollback(tag string) (int, error)
}

// Transaction represents a changelog transaction.
type Transaction interface {
	// Commit should attempt to commit the changes to to the changelog or
	// return an error.
	Commit() error
	// Rollback should undo changes to the changelog after a failed commit.
	Rollback() error
	// Exec should prepare to make a change to the changelog.
	Exec(query string) error
}
