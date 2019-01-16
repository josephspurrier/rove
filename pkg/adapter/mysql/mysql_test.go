package mysql_test

import (
	"testing"
	"time"

	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	c := mysql.Connection{}
	c.Username = "root"
	c.Password = "password"
	c.Hostname = "localhost"
	c.Port = 3306
	c.Database = "test"
	c.Parameter = "collation=utf8mb4_unicode_ci"

	// Test with database.
	dsn := c.DSN(true)
	assert.Equal(t, "root:password@tcp(localhost:3306)/test?collation=utf8mb4_unicode_ci", dsn)

	// Test without database.
	dsn = c.DSN(false)
	assert.Equal(t, "root:password@tcp(localhost:3306)/?collation=utf8mb4_unicode_ci", dsn)

	// Test without database and parameters.
	c.Parameter = ""
	dsn = c.DSN(false)
	assert.Equal(t, "root:password@tcp(localhost:3306)/", dsn)
}

func TestErrors(t *testing.T) {
	rr := new(mysql.MySQL)
	for _, v := range []error{
		rr.Initialize(),
		func() error {
			_, err := rr.ChangesetApplied("", "", "")
			return err
		}(),
		func() error {
			_, err := rr.BeginTx()
			return err
		}(),
		func() error {
			_, err := rr.Count()
			return err
		}(),
		rr.Insert("", "", "", time.Now(), 0, "", "", ""),
		func() error {
			_, err := rr.Changesets(false)
			return err
		}(),
		rr.Delete("", "", ""),
		rr.Tag("", "", "", ""),
		func() error {
			_, err := rr.Rollback("")
			return err
		}(),
	} {
		assert.Equal(t, mysql.ErrChangelogFailure, v)
	}
}
