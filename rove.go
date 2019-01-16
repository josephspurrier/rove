package rove

const (
	appVersion = "1.0"
)

// Rove contains the database migration information.
type Rove struct {
	// Verbose is whether information is written to the screen or not.
	Verbose bool
	// Checksum determines how operations continue if checksums don't match.
	Checksum ChecksumMode

	// file is the full path to the migration file.
	file string
	// changeset is text with changesets.
	changeset string
	// db is a migration.
	db Changelog
}

// ChecksumMode represents how to handle checksums on migrations.
type ChecksumMode int

const (
	// ChecksumThrowError [default] will throw an error if the checksum
	// in the changelog doesn't match the checksum of the changeset.
	ChecksumThrowError ChecksumMode = iota
	// ChecksumIgnore continues to process if the checksum
	// in the changelog doesn't match the checksum of the changeset.
	ChecksumIgnore
	// ChecksumUpdate updates the checksum in the changelog if the checksum
	// in the changelog doesn't match the checksum of the changeset.
	ChecksumUpdate
)

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
