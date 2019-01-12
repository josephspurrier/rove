package jsonfile

import (
	"fmt"
	"io/ioutil"
)

// Transaction is a file transaction.
type Transaction struct {
	data     string
	filename string
}

// NewTransaction returns a new file transaction.
func NewTransaction(filename string) *Transaction {
	t := new(Transaction)
	t.filename = filename
	return t
}

// Commit writes the data to the file.
func (t *Transaction) Commit() error {
	return ioutil.WriteFile(t.filename, []byte(t.data), 0644)
}

// Rollback reverts any changes due to a failed commit.
func (t *Transaction) Rollback() error {
	return nil
}

// Exec prepares the data write.
func (t *Transaction) Exec(query string, args ...interface{}) error {
	t.data = fmt.Sprintf(query, args...)
	return nil
}
