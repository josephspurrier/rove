package mysql

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/env"

	"github.com/jmoiron/sqlx"
)

const (
	// TestDatabaseName is the name of the test database.
	TestDatabaseName = "webapitest"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func setEnv(unique string) {
	os.Setenv(unique+"DB_HOSTNAME", "127.0.0.1")
	os.Setenv(unique+"DB_PORT", "3306")
	os.Setenv(unique+"DB_USERNAME", "root")
	os.Setenv(unique+"DB_PASSWORD", "")
	os.Setenv(unique+"DB_DATABASE", TestDatabaseName+unique)
	os.Setenv(unique+"DB_PARAMETER", "parseTime=true&allowNativePasswords=true&multiStatements=true")
}

func unsetEnv(unique string) {
	os.Unsetenv(unique + "DB_HOSTNAME")
	os.Unsetenv(unique + "DB_PORT")
	os.Unsetenv(unique + "DB_USERNAME")
	os.Unsetenv(unique + "DB_PASSWORD")
	os.Unsetenv(unique + "DB_DATABASE")
	os.Unsetenv(unique + "DB_PARAMETER")
}

// connectDatabase returns a test database connection.
func connectDatabase(dbSpecificDB bool, unique string) *sqlx.DB {
	dbc := new(Connection)
	err := env.Unmarshal(dbc, unique)
	if err != nil {
		fmt.Println("DB ENV Error:", err)
	}

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
	setEnv(unique)

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

	unsetEnv(unique)
}

// LoadDatabaseFromFile will set up the DB for the tests.
func LoadDatabaseFromFile(file string, usePrefix bool) (*sqlx.DB, string) {
	unique := ""
	var db *sqlx.DB

	var r *rove.Rove

	if usePrefix {
		db, unique = SetupDatabase()
		// Create a new MySQL database object.
		m := new(MySQL)
		m.DB = db
		r = rove.NewFileMigration(m, file)
	} else {
		m := new(MySQL)
		m.DB = db
		r = rove.NewFileMigration(m, file)
		setEnv(unique)
		db = connectDatabase(false, unique)
		_, err := db.Exec(`DROP DATABASE IF EXISTS ` + TestDatabaseName)
		if err != nil {
			fmt.Println("DB DROP SETUP Error:", err)
		}
		_, err = db.Exec(`CREATE DATABASE ` + TestDatabaseName + ` DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci`)
		if err != nil {
			fmt.Println("DB CREATE Error:", err)
		}

		db = connectDatabase(true, unique)
	}

	err := r.Migrate(0)
	if err != nil {
		log.Println("DB Error:", err)
	}

	return db, unique
}
