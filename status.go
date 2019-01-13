package rove

import (
	"fmt"

	"github.com/josephspurrier/rove/pkg/changeset"
)

// Status will output all changesets from the database table and will return
// an error, the last changeset, or a blank string. Check in that order.
func (r *Rove) Status() (*changeset.Record, error) {
	var last *changeset.Record

	// Get an array of changesets from the database.
	results, err := r.db.Changesets(false)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		if r.Verbose {
			fmt.Println("No changesets applied to the database.")
		}
		return nil, nil
	}

	if r.Verbose {
		fmt.Println("Changesets applied:")
	}

	// Loop through each changeset.
	for _, rs := range results {
		last = &rs
		if r.Verbose {
			fmt.Printf("%v\n", rs.String())
		}
	}

	return last, nil
}
