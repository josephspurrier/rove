package rove

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/josephspurrier/rove/pkg/changeset"
)

const (
	elementChangeset   = "--changeset "
	elementRollback    = "--rollback "
	elementInclude     = "--include "
	elementDescription = "--description "
	elementMemory      = "memory"
)

var (
	// ErrInvalidFormat is when a changeset is not found.
	ErrInvalidFormat = errors.New("invalid changeset format")
)

// parseFileToArray will parse a file into changesets.
func parseFileToArray(filename string) ([]changeset.Record, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseToArray(f, filename)
}

// parseToArray will split the migration into an ordered array.
func parseToArray(r io.Reader, filename string) ([]changeset.Record, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	// Array of changesets.
	arr := make([]changeset.Record, 0)

	for scanner.Scan() {
		// Get the line without leading or trailing spaces.
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines or liquibase header.
		if len(line) == 0 || strings.HasPrefix(line, "--liquibase") {
			continue
		}

		// Determine if the line is an `include`.
		if strings.HasPrefix(line, elementInclude) {
			// Load the file and add to the array.
			fp := strings.TrimPrefix(line, elementInclude)
			rfp := filepath.Join(filepath.Dir(filename), fp)
			cs, err := parseFileToArray(rfp)
			if err != nil {
				return nil, err
			}
			arr = append(arr, cs...)
			continue
		}

		// Start recording the changeset.
		if strings.HasPrefix(line, elementChangeset) {
			// Create a new changeset.
			cs := new(changeset.Record)
			cs.ParseHeader(strings.TrimPrefix(line, elementChangeset))
			cs.SetFileInfo(path.Base(filename), appVersion)
			arr = append(arr, *cs)
			continue
		}

		// If the length of the array is 0, then the first changeset is missing.
		if len(arr) == 0 {
			return nil, ErrInvalidFormat
		}

		// Determine if the line is a rollback.
		if strings.HasPrefix(line, elementRollback) {
			arr[len(arr)-1].AddRollback(strings.TrimPrefix(line, elementRollback))
			continue
		}

		// Determine if the line is a description.
		if strings.HasPrefix(line, elementDescription) {
			arr[len(arr)-1].AddDescription(strings.TrimPrefix(line, elementDescription))
			continue
		}

		// Determine if the line is comment, ignore it.
		if strings.HasPrefix(line, "--") {
			continue
		}

		// Add the line as a changeset.
		arr[len(arr)-1].AddChange(line)
	}

	// Perform a verification check on duplicates.
	_, err := parseArrayToMap(arr)

	return arr, err
}

// parseFileToMap will parse a file into a map.
func parseFileToMap(filename string) (map[string]changeset.Record, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseReaderToMap(f, filename)
}

// parseReaderToMap will parse a reader to a map.
func parseReaderToMap(r io.Reader, filename string) (map[string]changeset.Record, error) {
	arr, err := parseToArray(r, filename)
	if err != nil {
		return nil, err
	}

	return parseArrayToMap(arr)
}

// parseArrayToMap will convert an array of changesets to a map of changesets.
func parseArrayToMap(arr []changeset.Record) (map[string]changeset.Record, error) {
	m := make(map[string]changeset.Record)

	for _, cs := range arr {
		id := fmt.Sprintf("%v:%v:%v", cs.Author, cs.ID, cs.Filename)
		if _, found := m[id]; found {
			return nil, errors.New("duplicate entry found: " + id)
		}

		m[id] = cs
	}

	return m, nil
}

// loadChangesets will get the changesets based on the type of migration
// specified during the creation of the Rove object.
func (r *Rove) loadChangesets() (map[string]changeset.Record, error) {
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
	arr, err := parseReaderToMap(strings.NewReader(r.changeset), elementMemory)
	if err != nil {
		return nil, err
	}

	return arr, nil
}
