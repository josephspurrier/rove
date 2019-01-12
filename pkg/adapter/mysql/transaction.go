package mysql

import "database/sql"

// Tx is a database transaction.
type Tx struct {
	db *sql.Tx
}

// NewTx creates a new database transaction.
func NewTx(tx *sql.Tx) *Tx {
	return &Tx{
		db: tx,
	}
}

// Commit will commit changes to the database or return an error.
func (t *Tx) Commit() error {
	return t.db.Commit()
}

// Rollback will rollback changes to the database or return an error.
func (t *Tx) Rollback() error {
	return t.db.Rollback()
}

// Exec will run a query on the database.
func (t *Tx) Exec(query string) error {
	_, err := t.db.Exec(query)
	return err
}
