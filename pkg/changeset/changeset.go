// Package changeset handles operations on the text of a changeset.
package changeset

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidHeader is when the changeset header is invalid.
	ErrInvalidHeader = errors.New("invalid changeset header")
)

// Record is a changeset.
type Record struct {
	ID            string
	Author        string
	Filename      string
	MD5           string
	Description   string
	Version       string
	OrderExecuted int //FIXME: I don't think these should be here.

	change   []string
	rollback []string
}

// ParseHeader will parse the header information.
func (cs *Record) ParseHeader(line string) error {
	arr := strings.Split(line, ":")
	if len(arr) != 2 {
		return ErrInvalidHeader
	}

	cs.Author = arr[0]
	cs.ID = arr[1]

	return nil
}

// SetFileInfo will set the file information.
func (cs *Record) SetFileInfo(filename string, description string, version string) {
	cs.Filename = filename
	cs.Description = description
	cs.Version = version
}

// AddRollback will add a rollback command.
func (cs *Record) AddRollback(line string) {
	if len(cs.rollback) == 0 {
		cs.rollback = make([]string, 0)
	}
	cs.rollback = append(cs.rollback, line)
}

// AddChange will add a change command.
func (cs *Record) AddChange(line string) {
	if len(cs.change) == 0 {
		cs.change = make([]string, 0)
	}
	cs.change = append(cs.change, line)
}

// Changes will return all the changes.
func (cs *Record) Changes() string {
	return strings.Join(cs.change, "\n")
}

// Rollbacks will return all the rollbacks.
func (cs *Record) Rollbacks() string {
	return strings.Join(cs.rollback, "\n")
}

// Checksum returns an MD5 checksum for the changeset.
func (cs *Record) Checksum() string {
	return md5sum(cs.Changes())
}
