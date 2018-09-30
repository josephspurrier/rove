package rove

import (
	"database/sql"
	"fmt"
	"strings"
)

// Migrate will perform all the migrations in a file. If max is 0, all
// migrations are run.
func (r *Rove) Migrate(max int) error {
	// Create a changelog table.
	err := r.m.CreateChangelogTable()
	if err != nil {
		return err
	}

	// Get the changesets.
	arr, err := parseFileToArray(r.MigrationFile)
	if err != nil {
		return err
	}

	maxCounter := 0

	// Loop through each changeset.
	for _, cs := range arr {
		checksum := ""
		newChecksum := cs.Checksum()

		// Determine if the changeset was already applied.
		// Count the number of rows.
		checksum, err = r.m.ChangesetApplied(cs.id, cs.author, cs.filename)
		if err == nil {
			// Determine if the checksums match.
			if checksum != newChecksum {
				return fmt.Errorf("checksum does not match - existing changeset %v:%v has checksum %v, but new changeset has checksum %v",
					cs.author, cs.id, checksum, newChecksum)
			}

			if r.Verbose {
				fmt.Printf("Changeset already applied: %v:%v\n", cs.author, cs.id)
			}
			continue
		} else if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("internal error on changeset %v:%v - %v", cs.author, cs.id, err.Error())
		}

		arrQueries := strings.Split(cs.Changes(), ";")

		tx, err := r.m.BeginTx()
		if err != nil {
			return fmt.Errorf("sql error begin transaction - %v", err.Error())
		}

		// Loop through each change.
		for _, q := range arrQueries {
			if len(q) == 0 {
				continue
			}

			// Execute the query.
			err = tx.Exec(q)
			if err != nil {
				return fmt.Errorf("sql error on changeset %v:%v - %v", cs.author, cs.id, err.Error())
			}
		}

		err = tx.Commit()
		if err != nil {
			errr := tx.Rollback()
			if errr != nil {
				return fmt.Errorf("sql error on commit rollback %v:%v - %v", cs.author, cs.id, errr.Error())
			}
			return fmt.Errorf("sql error on commit %v:%v - %v", cs.author, cs.id, err.Error())
		}

		// Count the number of rows.
		count, err := r.m.Count()
		if err != nil {
			return err
		}

		// Insert the record.
		err = r.m.Insert(cs.id, cs.author, cs.filename, count+1, newChecksum,
			cs.description, cs.version)
		if err != nil {
			return err
		}

		if r.Verbose {
			fmt.Printf("Changeset applied: %v:%v\n", cs.author, cs.id)
		}

		// Only perform the maxium number of changes based on the max value.
		maxCounter++
		if max != 0 {
			if maxCounter >= max {
				break
			}
		}
	}

	return nil
}
