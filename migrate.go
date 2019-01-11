package rove

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/josephspurrier/rove/pkg/changeset"
)

// Migrate will perform all the migrations in a file. If max is 0, all
// migrations are run.
func (r *Rove) Migrate(max int) error {
	// Create a changelog table.
	err := r.db.CreateChangelogTable()
	if err != nil {
		return err
	}

	arr := make([]changeset.Info, 0)
	// If a file is specified, use it to build the array.
	if len(r.file) > 0 {
		// Get the changesets.
		arr, err = parseFileToArray(r.file)
		if err != nil {
			return err
		}
	} else {
		// Else use the changeset that was passed in.
		arr, err = parseToArray(strings.NewReader(r.changeset), elementMemory)
		if err != nil {
			return err
		}
	}

	maxCounter := 0

	// Loop through each changeset.
	for _, cs := range arr {
		checksum := ""
		newChecksum := cs.Checksum()

		// Determine if the changeset was already applied.
		// Count the number of rows.
		checksum, err = r.db.ChangesetApplied(cs.ID, cs.Author, cs.Filename)
		if err == nil {
			// Determine if the checksums match.
			if checksum != newChecksum {
				return fmt.Errorf("checksum does not match - existing changeset %v:%v has checksum %v, but new changeset has checksum %v",
					cs.Author, cs.ID, checksum, newChecksum)
			}

			if r.Verbose {
				fmt.Printf("Changeset already applied: %v:%v\n", cs.Author, cs.ID)
			}
			continue
		} else if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("internal error on changeset %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		arrQueries := strings.Split(cs.Changes(), ";")

		tx, err := r.db.BeginTx()
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
				return fmt.Errorf("sql error on changeset %v:%v - %v", cs.Author, cs.ID, err.Error())
			}
		}

		err = tx.Commit()
		if err != nil {
			errr := tx.Rollback()
			if errr != nil {
				return fmt.Errorf("sql error on commit rollback %v:%v - %v", cs.Author, cs.ID, errr.Error())
			}
			return fmt.Errorf("sql error on commit %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		// Count the number of rows.
		count, err := r.db.Count()
		if err != nil {
			return err
		}

		// Insert the record.
		err = r.db.Insert(cs.ID, cs.Author, cs.Filename, count+1, newChecksum,
			cs.Description, cs.Version)
		if err != nil {
			return err
		}

		if r.Verbose {
			fmt.Printf("Changeset applied: %v:%v\n", cs.Author, cs.ID)
		}

		// Only perform the maxium number of changes based on the max value.
		maxCounter++
		if max != 0 && maxCounter >= max {
			break
		}
	}

	return nil
}
