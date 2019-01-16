package mysql

import (
	"fmt"

	"github.com/josephspurrier/rove/pkg/env"

	_ "github.com/go-sql-driver/mysql" // MySQL driver.
	"github.com/jmoiron/sqlx"
)

// Connection holds the details for the MySQL connection.
type Connection struct {
	Hostname  string `json:"Hostname" env:"DB_HOSTNAME"`
	Port      int    `json:"Port" env:"DB_PORT"`
	Username  string `json:"Username" env:"DB_USERNAME"`
	Password  string `json:"Password" env:"DB_PASSWORD"`
	Name      string `json:"Database" env:"DB_NAME"`
	Parameter string `json:"Parameter" env:"DB_PARAMETER"`
}

// NewConnection returns the info required to make a connection to a MySQL
// database from environment variables. The optional prefix is used when reading
// environment variables.
func NewConnection(prefix string) (*Connection, error) {
	dbc := new(Connection)

	// Load the struct from environment variables.
	err := env.Unmarshal(dbc, prefix)
	if err != nil {
		return nil, err
	}

	return dbc, nil
}

// Connect to the database.
func (c Connection) Connect(includeDatabase bool) (*sqlx.DB, error) {
	// Connect to database and verify with a ping.
	return sqlx.Connect("mysql", c.DSN(includeDatabase))
}

// DSN returns the Data Source Name.
func (c Connection) DSN(includeDatabase bool) string {
	// Get parameters.
	param := c.Parameter

	// If parameter is specified, add a question mark. Users should not prefix
	// their parameter strings with a question mark.
	if len(c.Parameter) > 0 {
		param = "?" + c.Parameter
	}

	// Example:
	// root:password@tcp(localhost:3306)/test?collation=utf8mb4_unicode_ci
	s := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v", c.Username, c.Password,
		c.Hostname, c.Port, param)

	if includeDatabase {
		// Example:
		// root:password@tcp(localhost:3306)/?collation=utf8mb4_unicode_ci
		s = fmt.Sprintf("%v:%v@tcp(%v:%d)/%v%v", c.Username, c.Password,
			c.Hostname, c.Port, c.Name, param)
	}

	return s
}
