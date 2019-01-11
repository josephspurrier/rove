package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	c := Connection{}
	c.Username = "root"
	c.Password = "password"
	c.Hostname = "localhost"
	c.Port = 3306
	c.Database = "test"
	c.Parameter = "collation=utf8mb4_unicode_ci"

	// Test with database.
	dsn := c.dsn(true)
	assert.Equal(t, "root:password@tcp(localhost:3306)/test?collation=utf8mb4_unicode_ci", dsn)

	// Test without database.
	dsn = c.dsn(false)
	assert.Equal(t, "root:password@tcp(localhost:3306)/?collation=utf8mb4_unicode_ci", dsn)

	// Test without database and parameters.
	c.Parameter = ""
	dsn = c.dsn(false)
	assert.Equal(t, "root:password@tcp(localhost:3306)/", dsn)
}
