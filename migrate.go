package rove

import (
	"fmt"
	"strings"

	"github.com/josephspurrier/rove/pkg/changeset"
)

// Migrate will perform all the migrations in a file. If max is 0, all
// migrations are run.
func (r *Rove) Migrate(max int) error {
	// Create the object to store the changeset log.
	err := r.db.Initialize()
	if err != nil {
		return fmt.Errorf("error on changelog creation: %v", err)
	}

	arr := make([]changeset.Record, 0)
	// If a file is specified, use it to build the array.
	if len(r.file) > 0 {
		// Get the changesets.
		arr, err = parseFileToArray(r.file)
		if err != nil {
			return fmt.Errorf("error parsing file: %v", err)
		}
	} else {
		// Else use the changeset that was passed in.
		arr, err = parseToArray(strings.NewReader(r.changeset), elementMemory)
		if err != nil {
			return fmt.Errorf("error on parsing string: %v", err)
		}
	}

	if r.Verbose {
		fmt.Printf("Changesets applied (request: %v):\n", max)
	}

	maxCounter := 0

	// Loop through each changeset.
	for _, cs := range arr {
		newChecksum := cs.GenerateChecksum()

		// Determine if the changeset was already applied.
		// Count the number of rows.
		record, err := r.db.ChangesetApplied(cs.ID, cs.Author, cs.Filename)
		if err != nil {
			return fmt.Errorf("internal error on changeset %v:%v - %v", cs.Author, cs.ID, err.Error())
		} else if record != nil {
			// Determine if the checksums match.
			if record.Checksum != newChecksum {
				return fmt.Errorf("checksum does not match - existing changeset %v:%v has checksum %v, but new changeset has checksum %v",
					cs.Author, cs.ID, record.Checksum, newChecksum)
			}

			if r.Verbose {
				fmt.Printf("Already applied: %v\n", record.String())
			}
			continue
		}

		tx, err := r.db.BeginTx()
		if err != nil {
			return fmt.Errorf("error on begin transaction - %v", err.Error())
		}

		// Execute the query.
		err = tx.Exec(cs.Changes())
		if err != nil {
			return fmt.Errorf("error on changeset %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			errr := tx.Rollback()
			if errr != nil {
				return fmt.Errorf("error on commit rollback %v:%v - %v", cs.Author, cs.ID, errr.Error())
			}
			return fmt.Errorf("error on commit %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		// Count the number of rows.
		count, err := r.db.Count()
		if err != nil {
			return fmt.Errorf("error on counting changelog rows: %v", err)
		}

		// Insert the record.
		err = r.db.Insert(cs.ID, cs.Author, cs.Filename, count+1, newChecksum,
			cs.Description, cs.Version)
		if err != nil {
			return fmt.Errorf("error on inserting changelog record: %v", err)
		}

		// Query back the record.
		newRecord, err := r.db.ChangesetApplied(cs.ID, cs.Author, cs.Filename)
		if err != nil {
			return fmt.Errorf("error on querying changelog record: %v", err)
		}

		if r.Verbose {
			fmt.Printf("Applied: %v\n", newRecord.String())
		}

		// Only perform the maxium number of changes based on the max value.
		maxCounter++
		if max != 0 && maxCounter >= max {
			break
		}
	}

	return nil
}
