package rove

import (
	"fmt"
)

// Rollback will rollback a number of changesets to a tag.
func (r *Rove) Rollback(tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("error - rollback tag cannot be empty")
	}

	// Get the number of max queries to run.
	max, err := r.db.Rollback(tag)
	if err != nil {
		return err
	}

	if r.Verbose {
		fmt.Printf("Found tag (%v), will rollback (%v) changeset(s)\n", tag, max)
	}

	// Rollback the changesets.
	err = r.Reset(max)

	if r.Verbose {
		fmt.Printf("Rollback complete\n")
	}

	return err
}
