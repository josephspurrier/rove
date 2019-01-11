package changeset

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidHeader is when the changeset header is invalid.
	ErrInvalidHeader = errors.New("invalid changeset header")
)

// Info is a database changeset.
type Info struct {
	ID          string
	Author      string
	Filename    string
	MD5         string
	Description string
	Version     string

	change   []string
	rollback []string
}

// ParseHeader will parse the header information.
func (cs *Info) ParseHeader(line string) error {
	arr := strings.Split(line, ":")
	if len(arr) != 2 {
		return ErrInvalidHeader
	}

	cs.Author = arr[0]
	cs.ID = arr[1]

	return nil
}

// SetFileInfo will set the file information.
func (cs *Info) SetFileInfo(filename string, description string, version string) {
	cs.Filename = filename
	cs.Description = description
	cs.Version = version
}

// AddRollback will add a rollback command.
func (cs *Info) AddRollback(line string) {
	if len(cs.rollback) == 0 {
		cs.rollback = make([]string, 0)
	}
	cs.rollback = append(cs.rollback, line)
}

// AddChange will add a change command.
func (cs *Info) AddChange(line string) {
	if len(cs.change) == 0 {
		cs.change = make([]string, 0)
	}
	cs.change = append(cs.change, line)
}

// Changes will return all the changes.
func (cs *Info) Changes() string {
	return strings.Join(cs.change, "\n")
}

// Rollbacks will return all the rollbacks.
func (cs *Info) Rollbacks() string {
	return strings.Join(cs.rollback, "\n")
}

// Checksum returns an MD5 checksum for the changeset.
func (cs *Info) Checksum() string {
	return md5sum(cs.Changes())
}
