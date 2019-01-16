package rove

import (
	"fmt"
	"strings"
	"time"

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

	var arr []changeset.Record

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
			comment := ""

			// Determine if the checksums match.
			if record.Checksum != newChecksum {
				if r.Checksum == ChecksumThrowError {
					return fmt.Errorf("checksum does not match - existing changeset %v:%v has checksum %v, but new changeset has checksum %v",
						cs.Author, cs.ID, record.Checksum, newChecksum)
				} else if r.Checksum == ChecksumIgnore {
					// Ignore.
					comment = fmt.Sprintf("Ignoring checksum (%v), should be (%v)\n", record.Checksum, newChecksum)
				} else if r.Checksum == ChecksumUpdate {
					// Update the checksum.
					err = r.db.Update(record.ID, record.Author, record.Filename,
						record.DateExecuted, record.OrderExecuted, newChecksum,
						record.Description, record.Version)
					if err != nil {
						return fmt.Errorf("internal error on updating changeset %v:%v - %v", cs.Author, cs.ID, err.Error())

					}
					comment = fmt.Sprintf("Updated checksum from (%v) to (%v)\n", record.Checksum, newChecksum)

				}
			}

			if r.Verbose {
				fmt.Printf("Already applied: %v\n", record.String())
				if len(comment) > 0 {
					fmt.Printf("%v", comment)
				}
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
		err = r.db.Insert(cs.ID, cs.Author, cs.Filename, time.Now(), count+1,
			newChecksum, cs.Description, cs.Version)
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

		// Only perform the maximum number of changes based on the max value.
		maxCounter++
		if max != 0 && maxCounter >= max {
			break
		}
	}

	return nil
}
