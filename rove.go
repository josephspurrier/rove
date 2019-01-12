package rove

const (
	appVersion = "1.0"
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
	db Changelog
}

// NewFileMigration returns a file migration object.
func NewFileMigration(db Changelog, filename string) *Rove {
	return &Rove{
		db:   db,
		file: filename,
	}
}

// NewChangesetMigration returns a changeset migration object.
func NewChangesetMigration(db Changelog, changeset string) *Rove {
	return &Rove{
		db:        db,
		changeset: changeset,
	}
}
