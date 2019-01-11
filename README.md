# Rove

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/rove)](https://goreportcard.com/report/github.com/josephspurrier/rove)
[![Build Status](https://travis-ci.org/josephspurrier/rove.svg)](https://travis-ci.org/josephspurrier/rove)
[![Coverage Status](https://coveralls.io/repos/github/josephspurrier/rove/badge.svg?branch=master&timestamp=20180923-01)](https://coveralls.io/github/josephspurrier/rove?branch=master)
[![GoDoc](https://godoc.org/github.com/josephspurrier/rove?status.svg)](https://godoc.org/github.com/josephspurrier/rove)

## MySQL Database Migration Tool Based on Liquibase

This project is based off Liquibase, the database migration tool. It uses a slimmed down version of the DATABASECHANGELOG database table to store the applied changesets. It only supports SQL (no XML) migrations currently. For the most part, you can take your existing SQL migration files and use them with this tool. You can also import this package into your application and apply changesets without having to store the migrations on physical files. This allows you manage DB migrations inside of your application so you can distribute single binary applications.

**Note:** Do not run this tool on a database that already has Liquibase migrations applied - they are not compatible because the checksums are calculated is different. The DATABASECHANGELOG database table which is used for storing the changesets is also different.

## Dependencies

```
gopkg.in/alecthomas/kingpin.v2
github.com/go-sql-driver/mysql
github.com/jmoiron/sqlx
github.com/stretchr/testify/assert
```

## Quick Start with Docker Compose

You can build a docker image from this repository and set it up along with a MySQL container using Docker Compose.

```bash
# Create a docker image.
docker build -t rove:latest .

# Launch MySQL and the Rove tool with Docker Compose.
docker-compose up

# The database should now have the sample migrations applied.

# Shutdown the containers.
docker-compose down
```

## Testing Migrations in MySQL Docker

Use the following commands to start a MySQL container with Docker:

```bash
# Start MySQL without a password.
docker run -d --name=mysql57 -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql:5.7
# or start MySQL with a password.
docker run -d --name=mysql57 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=somepassword mysql:5.7

# Create the database via docker exec.
docker exec mysql57 sh -c 'exec mysql -uroot -e "CREATE DATABASE IF NOT EXISTS webapi DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;"'
# Or create the database manually.
CREATE DATABASE webapi DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

# Build the CLI tool.
go install

# Apply the database migrations without a password.
DB_USERNAME=root DB_HOSTNAME=127.0.0.1 DB_PORT=3306 DB_DATABASE=webapi rove migrate all testdata/success.sql
# or apply the database migrations with a password.
DB_USERNAME=root DB_PASSWORD=somepassword DB_HOSTNAME=127.0.0.1 DB_PORT=3306 DB_DATABASE=webapi rove migrate all testdata/success.sql
```

## Rove Usage

The Rove package can be used as a standalone CLI application or you can import the package into your own code.

### Rove via CLI

The Rove CLI app uses [Kingpin](https://github.com/alecthomas/kingpin) to handle command-line parsing. To view the info on any of the commands, you can use `--help` as an argument.

Here are the commands available:

```
usage: rove [<flags>] <command> [<args> ...]

Performs database migration tasks.

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  migrate all <file>
    Apply all changesets to the database.

  migrate up <count> <file>
    Apply a specific number of changesets to the database.

  migrate reset <file>
    Run all rollbacks on the database.

  migrate down <count> <file>
    Apply a specific number of rollbacks to the database.

  migrate status
    Output the list of migrations already applied to the database.
```

#### Environment Variables

The following environment variables can be read by the CLI app:

```
DB_USERNAME - database username
DB_PASSWORD - database password
DB_HOSTNAME - IP or hostname of the database
DB_PORT - port of the database
DB_DATABASE - name of the database
DB_PARAMETER - parameters to append to the database connection string
```

A full list of MySQL parameters can be found [here](https://github.com/go-sql-driver/mysql#parameters).

#### Example Commands

Here are examples of the commands:

```bash
# Set the environment variables to connect to the database.
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_HOSTNAME=127.0.0.1
export DB_PORT=3306
export DB_DATABASE=webapi
export DB_PARAMETER="collation=utf8mb4_unicode_ci&parseTime=true"

# Apply all of the changes from the SQL file to the database.
rove migrate all testdata/changeset.sql
# Output:
# Changeset applied: josephspurrier:1
# Changeset applied: josephspurrier:2
# Changeset applied: josephspurrier:3

# Try to apply all the changes again.
rove migrate all testdata/changeset.sql
# Output:
# Changeset already applied: josephspurrier:1
# Changeset already applied: josephspurrier:2
# Changeset already applied: josephspurrier:3

# Rollback all of the changes to the database.
rove migrate reset testdata/changeset.sql
# Output:
# Rollback applied: josephspurrier:3
# Rollback applied: josephspurrier:2
# Rollback applied: josephspurrier:1

# Apply only 1 new change to the database.
rove migrate up 1 testdata/success.sql
# Output:
# Changeset applied: josephspurrier:1

# Apply 1 more change to the database.
rove migrate up 1 testdata/changeset.sql
# Output:
# Changeset already applied: josephspurrier:1
# Changeset applied: josephspurrier:2

# Rollback only 1 change to the database.
rove migrate down 1 testdata/changeset.sql
# Output:
# Rollback applied: josephspurrier:2

# Show a list of migrations already applied to the database.
rove migrate status
# Output:
# Changeset: josephspurrier:1
```

### Rove via Package Import

Below is an example of how to include Rove in your own packages.

```go
var changesets = `
--changeset josephspurrier:1
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user_status;

--changeset josephspurrier:2
INSERT INTO user_status (id, status, created_at, updated_at, deleted) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);
--rollback TRUNCATE TABLE user_status;`

// Create a new MySQL database object.
db, err := database.NewMySQL(&database.Connection{
  Hostname:  "127.0.0.1",
  Username:  "root",
  Password:  "password",
  Database:  "main",
  Port:      3306,
  Parameter: "collation=utf8mb4_unicode_ci&parseTime=true",
})
if err != nil {
  log.Fatalln(err)
}

// Perform all migrations against the database.
r := rove.NewChangesetMigration(db, changesets)
r.Verbose = true
err = r.Migrate(0)
```