package testutil

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/josephspurrier/rove/pkg/adapter/mysql"

	"github.com/jmoiron/sqlx"
)

const (
	// TestDatabaseName is the name of the test database.
	TestDatabaseName = "webapitest"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SetEnv will set the environment variables.
func SetEnv(unique string) {
	os.Setenv(unique+"DB_HOSTNAME", "127.0.0.1")
	os.Setenv(unique+"DB_PORT", "3306")
	os.Setenv(unique+"DB_USERNAME", "root")
	os.Setenv(unique+"DB_PASSWORD", "")
	os.Setenv(unique+"DB_NAME", TestDatabaseName+unique)
	os.Setenv(unique+"DB_PARAMETER", "parseTime=true&allowNativePasswords=true&multiStatements=true")
}

// UnsetEnv will unset the environment variables.
func UnsetEnv(unique string) {
	os.Unsetenv(unique + "DB_HOSTNAME")
	os.Unsetenv(unique + "DB_PORT")
	os.Unsetenv(unique + "DB_USERNAME")
	os.Unsetenv(unique + "DB_PASSWORD")
	os.Unsetenv(unique + "DB_NAME")
	os.Unsetenv(unique + "DB_PARAMETER")
}

// Connection returns the test connection.
func Connection(unique string) *mysql.Connection {
	return &mysql.Connection{
		Hostname:  "127.0.0.1",
		Port:      3306,
		Username:  "root",
		Password:  "",
		Name:      TestDatabaseName + unique,
		Parameter: "parseTime=true&allowNativePasswords=true&multiStatements=true",
	}
}

// connectDatabase returns a test database connection.
func connectDatabase(dbSpecificDB bool, unique string) *sqlx.DB {
	dbc := Connection(unique)

	connection, err := dbc.Connect(dbSpecificDB)
	if err != nil {
		fmt.Println("DB Error:", err)
	}

	return connection
}

// SetupDatabase will create the test database and set the environment
// variables.
func SetupDatabase() (*sqlx.DB, string) {
	unique := "T" + fmt.Sprint(rand.Intn(500))

	db := connectDatabase(false, unique)
	_, err := db.Exec(`DROP DATABASE IF EXISTS ` + TestDatabaseName + unique)
	if err != nil {
		fmt.Println("DB DROP SETUP Error:", err)
	}
	_, err = db.Exec(`CREATE DATABASE ` + TestDatabaseName + unique + ` DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci`)
	if err != nil {
		fmt.Println("DB CREATE Error:", err)
	}

	return connectDatabase(true, unique), unique
}

// TeardownDatabase will destroy the test database and unset the environment
// variables.
func TeardownDatabase(unique string) {
	db := connectDatabase(false, unique)
	_, err := db.Exec(`DROP DATABASE IF EXISTS ` + TestDatabaseName + unique)
	if err != nil {
		fmt.Println("DB DROP TEARDOWN Error:", err)
	}
}
