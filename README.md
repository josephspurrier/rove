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
DB_USERNAME=root DB_HOSTNAME=127.0.0.1 DB_PORT=3306 DB_DATABASE=webapi rove all testdata/success.sql
# or apply the database migrations with a password.
DB_USERNAME=root DB_PASSWORD=somepassword DB_HOSTNAME=127.0.0.1 DB_PORT=3306 DB_DATABASE=webapi rove all testdata/success.sql
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

  all <file>
    Apply all changesets to the database.

  up <count> <file>
    Apply a specific number of changesets to the database.

  reset <file>
    Run all rollbacks on the database.

  down <count> <file>
    Apply a specific number of rollbacks to the database.

  status
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
rove all testdata/changeset.sql
# Output:
# Changeset applied: josephspurrier:1
# Changeset applied: josephspurrier:2
# Changeset applied: josephspurrier:3

# Try to apply all the changes again.
rove all testdata/changeset.sql
# Output:
# Changeset already applied: josephspurrier:1
# Changeset already applied: josephspurrier:2
# Changeset already applied: josephspurrier:3

# Rollback all of the changes to the database.
rove reset testdata/changeset.sql
# Output:
# Rollback applied: josephspurrier:3
# Rollback applied: josephspurrier:2
# Rollback applied: josephspurrier:1

# Apply only 1 new change to the database.
rove up 1 testdata/success.sql
# Output:
# Changeset applied: josephspurrier:1

# Apply 1 more change to the database.
rove up 1 testdata/changeset.sql
# Output:
# Changeset already applied: josephspurrier:1
# Changeset applied: josephspurrier:2

# Rollback only 1 change to the database.
rove down 1 testdata/changeset.sql
# Output:
# Rollback applied: josephspurrier:2

# Show a list of migrations already applied to the database.
rove status
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

## Adapters

Rove is designed to be extensible via adapters. There are a few adapters built include in the standard package:

* mysql
* jsonfile

You may also create your own adapters - see the `interface.go` file for interfaces your adapters must satisfy.

### Best Practices

When creating an adapter, will need:

- Struct that satisfies the `rove.Changelog` interface.
- Struct that satisfies the `rove.Transaction` interface.
- Table or data structure called a `Changelog` to persistently track the changes made by the Rove.

You should store the following fields (at a minimum) in your changelog. This will ensure your adapter can utilize all of the features of Rove.

- id
- author
- filename
- dateexecuted
- orderexecuted
- checksum
- description
- tag
- version

These fields are not provided by Rove but you should track them in the changelog:

- dateexecuted
- orderexecuted

### Example Changelog

Your changelog should look contain the same data as this table:

| id | author | filename | dateexecuted | orderexecuted | checksum | description | tag | version |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | josephspurrier | success.sql | 2019-01-12 16:04:16 | 1 | f0685b... | Create the user_status table. | NULL | 1.0 |
| 2 | josephspurrier | success.sql | 2019-01-12 16:04:16 | 2 | 3f81b0... |  | NULL | 1.0 |
| 3 | josephspurrier | success.sql | 2019-01-12 16:04:16 | 3 | 57cc0b... |  | NULL | 1.0 |

## Migration File Specifications

There are a few components of a changeset:

- Header: must be prefixed by "--changeset " and must follow this format: `author:id` (single line, required)
- Body: valid sql text (multi-line, required)
- Description: must be prefixed by "--description " (multi-line, optional)
- Rollback: must be prefixed "--rollback "  (multi-line, optional)
- Include: must be prefixed by "--include " and must follow this format: `relativefilename.sql` (single line, optional)
- Comments: any other line that starts with "--" (multi-line, optional)

The number of lines or between changesets does not matter.

Example migration file:

```sql
--changeset josephspurrier:1
--description Create the user status table.
--description Ensure no auto value on zero.
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

--include anotherfile.sql

--changeset josephspurrier:2
INSERT INTO user_status (id, status, created_at, updated_at, deleted) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);
--rollback TRUNCATE TABLE user_status;
```

### Header

The header is the unique identifier for the changeset. A changeset is unique is all of these fields don't match another changeset: id, author, and filename. You can have a changeset with the same id and author in two different files.

### Body

The body must be valid single or multi-line SQL queries. You can separate queries by semi-colons, but you must also pass in this parameter to the database connection: `multiStatements=true`. The checksum is based on an MD5 of this value. Any changes once the query has been applied to a database will throw an error message.

### Description

The description provides information about the changeset. It will be added as a value in the changelog table.

### Rollback

The rollback should be SQL reverting the changes made by the queries from the body.

### Include

The include allows you to reference other changeset files to load. The filename should be a relative path.

### Comments

Any comments at the beginning of the lines are ignored. They do not count towards the checksum.