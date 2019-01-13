package rove

import (
	"errors"
	"fmt"
)

// Reset will remove all migrations. If max is 0, all rollbacks are run.
func (r *Rove) Reset(max int) error {
	// Get an array of changesets from the database.
	results, err := r.db.Changesets(true)
	if err != nil {
		return err
	}

	// Get the changesets.
	m, err := r.loadChangesets()
	if err != nil {
		return err
	}

	if len(results) == 0 {
		if r.Verbose {
			fmt.Println("No rollbacks to perform.")
		}
		return nil
	}

	if r.Verbose {
		fmt.Printf("Changesets rollback (request: %v):\n", max)
	}

	maxCounter := 0

	// Loop through each changeset.
	for _, rs := range results {
		id := fmt.Sprintf("%v:%v:%v", rs.Author, rs.ID, rs.Filename)

		cs, ok := m[id]
		if !ok {
			return errors.New("changeset is missing: " + id)
		}

		tx, err := r.db.BeginTx()
		if err != nil {
			return fmt.Errorf("error on begin transaction - %v", err.Error())
		}

		// Execute the query.
		err = tx.Exec(cs.Rollbacks())
		if err != nil {
			return fmt.Errorf("error on rollback %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			errr := tx.Rollback()
			if errr != nil {
				return fmt.Errorf("error on commit rollback %v:%v - %v", cs.Author, cs.ID, errr.Error())
			}
			return fmt.Errorf("error on commit %v:%v - %v", cs.Author, cs.ID, err.Error())
		}

		// Delete the record.
		err = r.db.Delete(cs.ID, cs.Author, cs.Filename)
		if err != nil {
			return err
		}

		if r.Verbose {
			fmt.Printf("Applied: %v\n", rs.String())
		}

		// Only perform the maxium number of changes based on the max value.
		maxCounter++
		if max != 0 && maxCounter >= max {
			break
		}
	}

	return nil
}
