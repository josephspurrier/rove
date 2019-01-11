package rove

import (
	"fmt"
)

// Status will output all changesets from the database table and will return
// an error, the last changeset, or a blank string. Check in that order.
func (r *Rove) Status() (string, error) {
	last := ""

	// Get an array of changesets from the database.
	results, err := r.db.Changesets(false)
	if err != nil {
		return last, err
	}

	if len(results) == 0 {
		if r.Verbose {
			fmt.Println("No changesets applied to the database.")
		}
		return last, nil
	}

	// Loop through each changeset.
	for _, rs := range results {
		last = fmt.Sprintf("%v:%v", rs.Author, rs.ID)
		if r.Verbose {
			fmt.Println("Changeset:", last)
		}
	}

	return last, nil
}
