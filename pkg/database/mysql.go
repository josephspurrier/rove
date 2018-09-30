package database

import (
	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/env"

	"github.com/jmoiron/sqlx"
)

const (
	sqlChangelog = `CREATE TABLE IF NOT EXISTS databasechangelog (
	id varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	author varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	filename varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	dateexecuted datetime NOT NULL,
	orderexecuted int(11) NOT NULL,
	md5sum varchar(35) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	tag varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	version varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`
)

// MySQL is a MySQL database connection.
type MySQL struct {
	DB *sqlx.DB
}

// NewMySQL connects to the database.
// prefix is the optional prefix used when reading environment variables.
func NewMySQL(prefix string) (m *MySQL, err error) {
	m = new(MySQL)
	// Connect to the database.
	m.DB, err = connect(prefix)
	return m, err
}

// connect will connect to the database.
func connect(prefix string) (*sqlx.DB, error) {
	dbc := new(Connection)

	// Load the struct from environment variables.
	err := env.Unmarshal(dbc, prefix)
	if err != nil {
		return nil, err
	}

	return dbc.Connect(true)
}

// CreateChangelogTable will create the changelog table and return an error.
func (m *MySQL) CreateChangelogTable() (err error) {
	// Create the DATABASECHANGELOG.
	_, err = m.DB.Exec(sqlChangelog)
	if err != nil {
		return err
	}

	return nil
}

// ChangesetApplied will return the checksum and the error if it's found.
func (m *MySQL) ChangesetApplied(id, author, filename string) (checksum string, err error) {
	err = m.DB.Get(&checksum, `SELECT md5sum
	FROM databasechangelog
	WHERE id = ?
	AND author = ?
	AND filename = ?`, id, author, filename)
	return checksum, err
}

// BeginTx starts a transaction.
func (m *MySQL) BeginTx() (rove.Transaction, error) {
	// Begin a transaction.
	t, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}

	tx := NewTx(t)
	return tx, nil
}

// Count returns the number of changesets in the database.
func (m *MySQL) Count() (count int, err error) {
	err = m.DB.Get(&count, `SELECT COUNT(*) FROM databasechangelog`)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Insert will insert a new record into the database.
func (m *MySQL) Insert(id, author, filename string, count int, checksum, description, version string) error {
	_, err := m.DB.Exec(`
	INSERT INTO databasechangelog
	(id,author,filename,dateexecuted,orderexecuted,md5sum,description,version)
	VALUES(?,?,?,NOW(),?,?,?,?)
	`, id, author, filename, count, checksum, description, version)
	return err
}

// Changesets returns a list of the changesets from the database.
func (m *MySQL) Changesets() ([]rove.DBChangeset, error) {
	results := make([]rove.DBChangeset, 0)
	err := m.DB.Select(&results, `
		SELECT id, author, filename, orderexecuted
		FROM databasechangelog
		ORDER BY orderexecuted DESC;`)
	return results, err
}

// Delete will delete a changeset from the database.
func (m *MySQL) Delete(id, author, filename string) error {
	// Delete the record.
	_, err := m.DB.Exec(`
			DELETE FROM databasechangelog
			WHERE id = ? AND author = ? AND filename = ?
			LIMIT 1
			`, id, author, filename)
	return err
}
