package rove

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// LBChangelog represents a Liquibase database table.
type LBChangelog struct {
	ID            string    `db:"ID"`
	Author        string    `db:"AUTHOR"`
	Filename      string    `db:"FILENAME"`
	DateExecuted  time.Time `db:"DATEEXECUTED"`
	OrderExecuted int       `db:"ORDEREXECUTED"`
	Tag           *string   `db:"TAG"`
	Version       string    `db:"LIQUIBASE"`
}

// Convert convert a Liquibase table to a Rove table.
func (r *Rove) Convert(db *sqlx.DB) error {
	// Create the object to store the changeset log.
	err := r.db.Initialize()
	if err != nil {
		return fmt.Errorf("error on changelog creation: %v", err)
	}

	// Get the changesets.
	m, err := r.loadChangesets()
	if err != nil {
		return err
	}

	results := make([]LBChangelog, 0)
	err = db.Select(&results, `
	SELECT ID, AUTHOR, FILENAME, DATEEXECUTED, ORDEREXECUTED, TAG, LIQUIBASE
	FROM DATABASECHANGELOG ORDER BY ORDEREXECUTED ASC;
	`)
	if err != nil {
		return err
	}

	log.Println(results)
	log.Println(m)

	// Loop through the Liquibase migrations.
	for _, cs := range results {
		// Determine if the changeset was already applied.
		record, err := r.db.ChangesetApplied(cs.ID, cs.Author, cs.Filename)
		if err != nil {
			return fmt.Errorf("internal error on changeset %v:%v - %v", cs.Author, cs.ID, err.Error())
		} else if record != nil {
			// Skip changesets that are already converted.
			continue
		}

		// Get the changeset from the map.
		id := fmt.Sprintf("%v:%v:%v", cs.Author, cs.ID, cs.Filename)
		newCS, ok := m[id]
		if !ok {
			return errors.New("changeset is missing: " + id)
		}

		// Insert the record.
		err = r.db.Insert(cs.ID, cs.Author, cs.Filename, cs.DateExecuted,
			cs.OrderExecuted, newCS.GenerateChecksum(),
			"", "liquibase "+cs.Version)
		if err != nil {
			return fmt.Errorf("error on inserting changelog record: %v", err)
		}

		// Query back the record.
		newRecord, err := r.db.ChangesetApplied(cs.ID, cs.Author, cs.Filename)
		if err != nil {
			return fmt.Errorf("error on querying changelog record: %v", err)
		}

		if r.Verbose {
			fmt.Printf("Converted: %v\n", newRecord.String())
		}
	}

	return nil
}
