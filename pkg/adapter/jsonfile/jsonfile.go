// Package jsonfile is a JSON file changelog adapter.
package jsonfile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/josephspurrier/rove/pkg/changeset"
)

// Info is a JSON changelog.
type Info struct {
	filename string
	db       *mysql.MySQL
}

// New sets the filename and returns an object that satisfies the rove.Changelog
// interface.
func New(filename string, db *mysql.MySQL) (m *Info, err error) {
	m = new(Info)
	m.filename = filename
	m.db = db
	return m, nil
}

// Initialize creates the changelog.
func (m *Info) Initialize() (err error) {
	// If the file doesn't exist, create the file.
	if _, err := os.Stat(m.filename); os.IsNotExist(err) {
		err = ioutil.WriteFile(m.filename, []byte("[]"), 0644)
	}

	return nil
}

// ChangesetApplied returns the checksum if it's found, an error if there was an
// issue, or a blank checksum with no error if it's not found.
func (m *Info) ChangesetApplied(id, author, filename string) (checksum string, err error) {
	// Read the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return "", err
	}

	// Convert to JSON.
	data := make([]changeset.Record, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return "", err
	}

	// Loop through to find the data.
	for _, cs := range data {
		// If found, return the checksum.
		if cs.ID == id && cs.Author == author && cs.Filename == filename {
			return cs.MD5, nil
		}
	}

	// If not found, return a blank string.
	return "", nil
}

// BeginTx starts a transaction.
func (m *Info) BeginTx() (rove.Transaction, error) {
	//tx := NewTransaction(m.filename)
	//return tx, nil
	return m.db.BeginTx()
}

// Count returns the number of changesets in the changelog.
func (m *Info) Count() (count int, err error) {
	// Read the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return 0, err
	}

	// Convert to JSON.
	data := make([]changeset.Record, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

// Insert will insert a new record into the database.
func (m *Info) Insert(id, author, filename string, count int, checksum, description, version string) error {
	// Read the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return err
	}

	// Convert to JSON.
	data := make([]changeset.Record, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	data = append(data, changeset.Record{
		Author:      author,
		Description: description,
		Filename:    filename,
		ID:          id,
		MD5:         checksum,
		Version:     version,
	})

	b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.filename, b, 0644)
}

// Changesets returns a list of the changesets from the database in ascending
// order (false) or descending order (true).
func (m *Info) Changesets(reverse bool) ([]rove.Changeset, error) {
	// Read the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return nil, err
	}

	// Convert to JSON.
	results := make([]changeset.Record, 0)
	err = json.Unmarshal(b, &results)
	if err != nil {
		return nil, err
	}

	// Copy from one struct to another.
	out := make([]rove.Changeset, 0)
	for _, i := range results {
		if reverse {
			out = append([]rove.Changeset{{
				Author:        i.Author,
				Filename:      i.Filename,
				ID:            i.ID,
				OrderExecuted: i.OrderExecuted,
			}}, out...)
		} else {
			out = append(out, rove.Changeset{
				Author:        i.Author,
				Filename:      i.Filename,
				ID:            i.ID,
				OrderExecuted: i.OrderExecuted,
			})
		}
	}

	return out, err
}

// Delete will delete a changeset from the database.
func (m *Info) Delete(id, author, filename string) error {
	// Read the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return err
	}

	// Convert to JSON.
	data := make([]changeset.Record, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	newData := make([]changeset.Record, 0)
	for _, cs := range data {
		if cs.ID == id && cs.Author == author && cs.Filename == filename {
			// skip
		} else {
			newData = append(newData, cs)
		}
	}

	b, err = json.Marshal(newData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.filename, b, 0644)
}
