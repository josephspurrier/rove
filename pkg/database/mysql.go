package database

import (
	"github.com/josephspurrier/rove"

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
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
)

// dbchangeset contains a single database record change.
type dbchangeset struct {
	ID            string `db:"id"`
	Author        string `db:"author"`
	Filename      string `db:"filename"`
	OrderExecuted int    `db:"orderexecuted"`
}

// MySQL is a MySQL database connection.
type MySQL struct {
	DB *sqlx.DB
}

// NewMySQL connects to the database and returns an object that satisfies the
// rove.Migration interface.
func NewMySQL(c *Connection) (m *MySQL, err error) {
	// Connect to the database.
	m = new(MySQL)
	m.DB, err = c.Connect(true)
	return m, err
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
	VALUES(?,?,?,NOW(),?,?,?,?)`,
		id, author, filename, count, checksum, description, version)
	return err
}

// Changesets returns a list of the changesets from the database in ascending
// order (false) or descending order (true).
func (m *MySQL) Changesets(reverse bool) ([]rove.DBChangeset, error) {
	order := "ASC"
	if reverse {
		order = "DESC"
	}

	results := make([]dbchangeset, 0)
	err := m.DB.Select(&results, `
	SELECT id, author, filename, orderexecuted
	FROM databasechangelog
	ORDER BY orderexecuted `+order)

	// Copy from one struct to another.
	out := make([]rove.DBChangeset, 0)
	for _, i := range results {
		out = append(out, rove.DBChangeset(i))
	}

	return out, err
}

// Delete will delete a changeset from the database.
func (m *MySQL) Delete(id, author, filename string) error {
	// Delete the record.
	_, err := m.DB.Exec(`
	DELETE FROM databasechangelog
	WHERE id = ? AND author = ? AND filename = ? LIMIT 1`, id, author, filename)
	return err
}
