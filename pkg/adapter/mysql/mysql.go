// Package mysql is a MySQL changelog adapter.
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/changeset"

	"github.com/jmoiron/sqlx"
)

const (
	tableName   = "rovechangelog"
	createQuery = `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
	id varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
	author varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
	filename varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
	dateexecuted datetime NOT NULL,
	orderexecuted int(11) NOT NULL,
	checksum char(32) COLLATE utf8mb4_unicode_ci NOT NULL,
	description varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
	tag varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL UNIQUE,
	version varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
)

var (
	// ErrChangelogFailure occurs when the connection is not set up properly.
	ErrChangelogFailure = errors.New("error with changelog setup")
	// ErrTransactionFuncMissing occurs when the transaction function is missing.
	ErrTransactionFuncMissing = errors.New("error transaction func is missing")
)

// dbchangeset contains a single database record change.
type dbchangeset struct {
	ID            string    `db:"id"`
	Author        string    `db:"author"`
	Filename      string    `db:"filename"`
	DateExecuted  time.Time `db:"dateexecuted"`
	OrderExecuted int       `db:"orderexecuted"`
	Checksum      string    `db:"checksum"`
	Description   string    `db:"description"`
	Tag           *string   `db:"tag"`
	Version       string    `db:"version"`
}

// MySQL is a MySQL database changelog.
type MySQL struct {
	DB              *sqlx.DB
	TableName       string
	InitializeQuery string
	TransactionFunc func(tx *sql.Tx) rove.Transaction
}

// New connects to the database and returns an object that satisfies the
// rove.Changelog interface.
func New(c *Connection) (m *MySQL, err error) {
	// Connect to the database.
	m = new(MySQL)
	m.DB, err = c.Connect(true)

	// Set the default table, create, and transaction.
	m.TableName = tableName
	m.InitializeQuery = createQuery
	m.TransactionFunc = func(tx *sql.Tx) rove.Transaction {
		return NewTx(tx)
	}

	return m, err
}

// Initialize will create the changelog table or return an error.
func (m *MySQL) Initialize() (err error) {
	if m.DB == nil {
		return ErrChangelogFailure
	}

	// Create the table.
	_, err = m.DB.Exec(m.InitializeQuery)
	if err != nil {
		return err
	}

	return nil
}

// ToRecord converts a dbchangeset to a changeset.Record.
func (m *MySQL) ToRecord(cs dbchangeset) *changeset.Record {
	tag := ""
	if cs.Tag != nil {
		tag = *cs.Tag
	}

	return &changeset.Record{
		ID:            cs.ID,
		Author:        cs.Author,
		Filename:      cs.Filename,
		DateExecuted:  cs.DateExecuted,
		OrderExecuted: cs.OrderExecuted,
		Checksum:      cs.Checksum,
		Description:   cs.Description,
		Tag:           tag,
		Version:       cs.Version,
	}
}

// ChangesetApplied returns the checksum from the database if it's found, an
// error if there was an issue, or nil with no error if it's not
// found.
func (m *MySQL) ChangesetApplied(id, author, filename string) (*changeset.Record, error) {
	if m.DB == nil {
		return nil, ErrChangelogFailure
	}

	var cs dbchangeset
	err := m.DB.Get(&cs, `
	SELECT * FROM `+m.TableName+`
	WHERE id = ?
	AND author = ?
	AND filename = ?`, id, author, filename)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return m.ToRecord(cs), err
}

// BeginTx starts a transaction.
func (m *MySQL) BeginTx() (rove.Transaction, error) {
	if m.DB == nil {
		return nil, ErrChangelogFailure
	}

	if m.TransactionFunc == nil {
		return nil, ErrTransactionFuncMissing
	}

	// Begin a transaction.
	t, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}

	return m.TransactionFunc(t), nil
}

// Count returns the number of changesets in the database.
func (m *MySQL) Count() (count int, err error) {
	if m.DB == nil {
		return 0, ErrChangelogFailure
	}

	err = m.DB.Get(&count, `SELECT COUNT(*) FROM `+m.TableName)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Insert will insert a new record into the database.
func (m *MySQL) Insert(id, author, filename string, dateexecuted time.Time,
	count int, checksum, description, version string) error {
	if m.DB == nil {
		return ErrChangelogFailure
	}

	_, err := m.DB.Exec(`
	INSERT INTO `+m.TableName+`
	(id,author,filename,dateexecuted,orderexecuted,checksum,description,version)
	VALUES(?,?,?,?,?,?,?,?)`,
		id, author, filename, dateexecuted, count, checksum, description, version)
	return err
}

// Update will update a record from the database.
func (m *MySQL) Update(id, author, filename string, dateexecuted time.Time,
	count int, checksum, description, version string) error {
	if m.DB == nil {
		return ErrChangelogFailure
	}

	_, err := m.DB.Exec(`
	UPDATE `+m.TableName+`
	SET
		dateexecuted = ?,
		orderexecuted = ?,
		checksum = ?,
		description = ?,
		version = ?
	WHERE
		id = ? AND 
		author = ? AND
		filename = ?
	LIMIT 1`,
		dateexecuted, count, checksum, description, version, id, author, filename)
	return err
}

// Changesets returns a list of the changesets from the database in ascending
// order (false) or descending order (true).
func (m *MySQL) Changesets(reverse bool) ([]changeset.Record, error) {
	if m.DB == nil {
		return nil, ErrChangelogFailure
	}

	order := "ASC"
	if reverse {
		order = "DESC"
	}

	results := make([]dbchangeset, 0)
	err := m.DB.Select(&results, `
	SELECT *
	FROM `+m.TableName+`
	ORDER BY orderexecuted `+order)

	// Copy from one struct to another.
	out := make([]changeset.Record, 0)
	for _, i := range results {
		out = append(out, *m.ToRecord(i))
	}

	return out, err
}

// Delete will delete a changeset from the database.
func (m *MySQL) Delete(id, author, filename string) error {
	if m.DB == nil {
		return ErrChangelogFailure
	}

	// Delete the record.
	_, err := m.DB.Exec(`
	DELETE FROM `+m.TableName+`
	WHERE id = ? AND author = ? AND filename = ? LIMIT 1`, id, author, filename)
	return err
}

// Tag will add a tag to the record.
func (m *MySQL) Tag(id, author, filename, tag string) error {
	if m.DB == nil {
		return ErrChangelogFailure
	}

	_, err := m.DB.Exec(`
	UPDATE `+m.TableName+`
	SET tag=?
	WHERE id = ? AND author = ? AND filename = ? LIMIT 1`,
		tag, id, author, filename)

	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}

	if me.Number == 1062 {
		return fmt.Errorf("tag already found in database: %v", tag)
	}

	return err
}

// Rollback return how many changesets to rollback.
func (m *MySQL) Rollback(tag string) (int, error) {
	if m.DB == nil {
		return 0, ErrChangelogFailure
	}

	count := 0
	err := m.DB.Get(&count, `
	SELECT count(*) FROM `+m.TableName+`
	WHERE id > (
		SELECT id FROM `+m.TableName+` WHERE tag = ?
	)`, tag)

	if count == 0 {
		return 0, fmt.Errorf("tag not found in database or no rollbacks to perform: %v", tag)
	}

	return count, err
}
