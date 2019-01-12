package jsonfile

import (
	"encoding/json"
	"io/ioutil"
)

// load returns an array from the JSON file.
func (m *Info) load() ([]Changeset, error) {
	// Load the file into memory.
	b, err := ioutil.ReadFile(m.filename)
	if err != nil {
		return nil, err
	}

	// Convert to JSON.
	data := make([]Changeset, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// save writes the array to a JSON file.
func (m *Info) save(data []Changeset) error {
	// Convert the data into JSON.
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the data to JSON.
	return ioutil.WriteFile(m.filename, b, m.FileMode)
}
