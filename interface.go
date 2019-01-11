package rove

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
	Changesets(reverse bool) ([]Change, error)
	// Delete the changeset from the database or return an error if there is a
	// problem running the query.
	Delete(id, author, filename string) error
}

// Change contains a single database record change.
type Change struct {
	ID            string
	Author        string
	Filename      string
	OrderExecuted int
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
