package rove

import (
	"errors"
	"fmt"
)

// Tag will tag the latest changelog with a string to allow for rollbacks.
func (r *Rove) Tag(tag string) error {
	// Get the changesets.
	m, err := r.loadChangesets()
	if err != nil {
		return err
	}

	// Get an array of changesets from the database.
	results, err := r.db.Changesets(true)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		if r.Verbose {
			fmt.Println("No changesets to tag.")
		}
		return nil
	}

	rs := results[0]
	id := fmt.Sprintf("%v:%v:%v", rs.Author, rs.ID, rs.Filename)

	// Get the changeset.
	_, ok := m[id]
	if !ok {
		return errors.New("changeset is missing: " + id)
	}

	// Tag the changeset.
	err = r.db.Tag(rs.ID, rs.Author, rs.Filename, tag)
	if err != nil {
		return fmt.Errorf("error on tag - %v", err.Error())
	}

	if r.Verbose {
		fmt.Printf("Tag applied: %v on %v:%v\n", tag, rs.Author, rs.ID)
	}

	return nil
}
