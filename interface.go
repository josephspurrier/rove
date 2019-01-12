package rove

// Changelog represents a list of operations for a changelog.
type Changelog interface {
	// Initialize should perform any work to set up the changelog or return an
	// error.
	Initialize() error
	// ChangesetApplied should return the checksum from the changelog of a
	// matching changeset, an error, or a blank string if the changeset doesn't
	// exist in the changelog.
	ChangesetApplied(id string, author string, filename string) (checksum string, err error)
	// BeginTx should start a transaction on the changelog.
	BeginTx() (Transaction, error)
	// Count should return the number of changesets in the changelog or return
	// an error.
	Count() (count int, err error)
	// Insert should add a new changeset to the changelog or return an error.
	Insert(id, author, filename string, count int, checksum, description, version string) error
	// Changesets should return a list of the changesets or return an error.
	Changesets(reverse bool) ([]Changeset, error)
	// Delete should remove the changeset from the changelog or return an error.
	Delete(id, author, filename string) error
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

// Changeset is a single changeset.
type Changeset struct {
	ID       string
	Author   string
	Filename string
}
