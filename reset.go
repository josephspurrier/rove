package rove

import (
	"errors"
	"fmt"
	"strings"
)

// loadChangesets will get the changesets based on the type of migration
// specified during the creation of the Rove object.
func (r *Rove) loadChangesets() (map[string]Changeset, error) {
	// Use the file to get the changesets first.
	if len(r.file) > 0 {
		// Get the changesets in a map.
		m, err := parseFileToMap(r.file)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	// Else use the changeset that was passed in.
	arr, err := parseReaderToMap(strings.NewReader(r.changeset), "memory")
	if err != nil {
		return nil, err
	}

	return arr, nil
}

// Reset will remove all migrations. If max is 0, all rollbacks are run.
func (r *Rove) Reset(max int) error {
	// Get the changesets.
	m, err := r.loadChangesets()
	if err != nil {
		return err
	}

	// Get an array of changesets from the database.
	results, err := r.db.Changesets()
	if err != nil {
		return err
	}

	if len(results) == 0 {
		if r.Verbose {
			fmt.Println("No rollbacks to perform.")
			return nil
		}
	}

	maxCounter := 0

	// Loop through each changeset.
	for _, rs := range results {
		id := fmt.Sprintf("%v:%v:%v", rs.Author, rs.ID, rs.Filename)

		cs, ok := m[id]
		if !ok {
			return errors.New("changeset is missing: " + id)
		}

		arrQueries := strings.Split(cs.Rollbacks(), ";")

		tx, err := r.db.BeginTx()
		if err != nil {
			return fmt.Errorf("sql error begin transaction - %v", err.Error())
		}

		// Loop through each rollback.
		for _, q := range arrQueries {
			if len(q) == 0 {
				continue
			}

			// Execute the query.
			err = tx.Exec(q)
			if err != nil {
				return fmt.Errorf("sql error on rollback %v:%v - %v", cs.author, cs.id, err.Error())
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

		// Delete the record.
		err = r.db.Delete(cs.id, cs.author, cs.filename)
		if err != nil {
			return err
		}

		if r.Verbose {
			fmt.Printf("Rollback applied: %v:%v\n", cs.author, cs.id)
		}

		// Only perform the maxium number of changes based on the max value.
		maxCounter++
		if max != 0 && maxCounter >= max {
			break
		}
	}

	return nil
}
