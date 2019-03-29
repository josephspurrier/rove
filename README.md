# Rove

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/rove)](https://goreportcard.com/report/github.com/josephspurrier/rove)
[![Build Status](https://travis-ci.org/josephspurrier/rove.svg)](https://travis-ci.org/josephspurrier/rove)
[![Coverage Status](https://coveralls.io/repos/github/josephspurrier/rove/badge.svg?branch=master&timestamp=20190115-01)](https://coveralls.io/github/josephspurrier/rove?branch=master)
[![GoDoc](https://godoc.org/github.com/josephspurrier/rove?status.svg)](https://godoc.org/github.com/josephspurrier/rove)

## MySQL Database Migration Tool Inspired by Liquibase

The primary motivation behind this tool is to provide a simple and quick Go (instead of Java) based database migration tool that allows loading migrations from anywhere, including from inside your code so you can distribute single binary applications. You write the migration and rollback SQL, Rove will apply it for you properly.

### How do migrations work?

Database migrations are necessary when you make changes to your applications. You may need a new table or a new column so you have to write the SQL commands to make the changes. The tricky piece is when you perform an upgrade, how do you manage which SQL queries will run? Do you run all of them again and then the new ones after? Or is there an easy way to track which queries have been run so you only run new ones? What if you have to rollback your database because of a feature that was released too early and is causing problem? How do you manage those queries? You can definitely write your own code to manage the migration process, but Rove makes the process much easier for you. You also don't have to convert your SQL code to a another format like JSON or XML, you can just add a few comments around it and Rove will handle the rest.

### How does Rove work?

You'll need to write your changes queries and rollback queries in migration files. These are plain SQL files that can be imported directly into MySQL. Rove just uses comments to help break them into smaller manageable pieces. When you run tell Rove to apply your changes, a table called `rovechangelog` is created in the database to track which changesets have been applied and metadata about them. The tool will ensure no changes have been made to the existing changesets that are already in the database. Changeset checksums are then compared against the changelog table checksums. Any new changesets that are not in the changelog are applied to the database and then a new record is inserted into the changelog for each changeset. Rove supports labeling changesets with a `tag` as well as rolling back to specific tags.

### Rove vs Liquibase

Rove and Liquibase use different changelog tables. Rove includes MySQL out of the box, but it supports adding your own adapters to work with any type of data storage. The Rove changesets can use a very similar plain SQL (no XML or JSON) file format for simplicity and portability. For the most teams, you'll be able use your existing SQL migration files with Rove without making any changes.

To assist with switching from Liquibase to Rove, you can use the CLI tool with the `rove convert` argument to convert a Liquibase `DATABASECHANGELOG` table to a Rove `rovechangelog` table. If you don't run the `rove convert` command first on a database that was originally managed by Liquibase, Rove will try to rerun the same migrations over again if you use the same migration files. The tools use different changelog table names, table schemas, and use different methods for calculating their checksums.

## Dependencies

These are the dependencies required to build Rove.

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
rove all testdata/success.sql --hostname=127.0.0.1 --port=3306 --username=root --name=webapi
# or apply the database migrations with a password.
rove all testdata/success.sql --hostname=127.0.0.1 --port=3306 --username=root --password=somepassword --name=webapi
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
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --checksum-mode=CHECKSUM-MODE  Set how to handle checksums that don't match [error (default), ignore, update].
  --hostname=HOSTNAME            Database hostname or IP [string].
  --port=PORT                    Database port [int].
  --username=USERNAME            Database username [string].
  --password=PASSWORD            Database password [string].
  --name=NAME                    Database name [string].
  --parameter=PARAMETER          Database parameters [string].
  --envprefix=ENVPREFIX          Prefix for environment variables.

Commands:
  help [<command>...]
    Show help.

  all <file>
    Apply all changesets to the database.

  up <count> <file>
    Apply a specific number of changesets to the database.

  reset <file>
    Apply all rollbacks to the database.

  down <count> <file>
    Apply a specific number of rollbacks to the database.

  tag <name> <file>
    Apply a tag to the latest changeset in the database.

  rollback <name> <file>
    Run all rollbacks until the specified tag on the database.

  convert <file>
    Convert a Liquibase changelog table to a Rove changelog table.

  status
    Output the list of changesets already applied to the database.
```

#### Database Connection Variables

You can either use the database flags or you can set the environment variables below to connect to the database. You can also prefix the environment variables using the `--envprefix` flag.

```
DB_USERNAME - database username
DB_PASSWORD - database password
DB_HOSTNAME - IP or hostname of the database
DB_PORT - port of the database
DB_NAME - name of the database
DB_PARAMETER - parameters to append to the database connection string
```

A full list of MySQL parameters can be found [here](https://github.com/go-sql-driver/mysql#parameters).

#### Example Commands and Output

These examples will show how to interact with Rove and what the output will look like.

```bash
# Set the environment variables to connect to the database.
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_HOSTNAME=127.0.0.1
export DB_PORT=3306
export DB_NAME=webapi
export DB_PARAMETER="collation=utf8mb4_unicode_ci&parseTime=true"

# Apply all of the changes from the SQL file to the database.
rove all testdata/changeset.sql
# Output:
# Changesets applied (request: 0):
# Applied: 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']
# Applied: 2) josephspurrier:2 (success.sql) e3065c58bff00322c73eab057427f557 [tag='']
# Applied: 3) josephspurrier:3 (success.sql) 57cc0b1c45cb72032bcaed07483d243d [tag='']

# Try to apply all the changes again.
rove all testdata/changeset.sql
# Output:
# Changesets applied (request: 0):
# Already applied: 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']
# Already applied: 2) josephspurrier:2 (success.sql) e3065c58bff00322c73eab057427f557 [tag='']
# Already applied: 3) josephspurrier:3 (success.sql) 57cc0b1c45cb72032bcaed07483d243d [tag='']

# Rollback all of the changes to the database.
rove reset testdata/changeset.sql
# Output:
# Changesets rollback (request: 0):
# Applied: 3) josephspurrier:3 (success.sql) 57cc0b1c45cb72032bcaed07483d243d [tag='']
# Applied: 2) josephspurrier:2 (success.sql) e3065c58bff00322c73eab057427f557 [tag='']
# Applied: 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']

# Apply only 1 new change to the database.
rove up 1 testdata/success.sql
# Output:
# Changesets applied (request: 1):
# Applied: 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']

# Apply 1 more change to the database.
rove up 1 testdata/changeset.sql
# Output:
# Changesets applied (request: 1):
# Already applied: 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']
# Applied: 2) josephspurrier:2 (success.sql) e3065c58bff00322c73eab057427f557 [tag='']

# Rollback only 1 change to the database.
rove down 1 testdata/changeset.sql
# Output:
# Changesets rollback (request: 1):
# Applied: 2) josephspurrier:2 (success.sql) e3065c58bff00322c73eab057427f557 [tag='']

# Show a list of migrations already applied to the database.
rove status
# Output:
# Changesets applied:
# 1) josephspurrier:1 (success.sql) b7a8d1c3ea1cc2dc28a1de0e23628250 [tag='']
```

### Rove via Package Import

Below is an example of how to include Rove in your own Go applications.

```go
var changesets = `
--changeset josephspurrier:1
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (id)
);
--rollback DROP TABLE user_status;

--changeset josephspurrier:2
INSERT INTO user_status (id, status, created_at, updated_at, deleted) VALUES
(1, 'active',   CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0),
(2, 'inactive', CURRENT_TIMESTAMP,  CURRENT_TIMESTAMP,  0);
--rollback TRUNCATE TABLE user_status;`

// Create a new MySQL database object.
db, err := mysql.New(&mysql.Connection{
  Hostname:  "127.0.0.1",
  Username:  "root",
  Password:  "password",
  Name:      "main",
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

Rove is designed to be extensible via adapters. There is one adapter included in the package:

* mysql

You may also create your own adapters - see the `interface.go` file for interfaces your adapters must satisfy.

### Best Practices

When creating an adapter, will need:

- Struct that satisfies the `rove.Changelog` interface.
- Struct that satisfies the `rove.Transaction` interface.
- Table or data structure to use as the `changelog` to persistently track the changes made by the Rove.

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

### Example Changelog

Your changelog should contain the same fields as this table:

| id  | author         | filename    | dateexecuted        | orderexecuted | checksum  | description                   | tag  | version |
| --- | -------------- | ----------- | ------------------- | ------------- | --------- | ----------------------------- | ---- | ------- |
| 1   | josephspurrier | success.sql | 2019-01-12 16:04:16 | 1             | f0685b... | Create the user_status table. | NULL | 1.0     |
| 2   | josephspurrier | success.sql | 2019-01-12 16:04:16 | 2             | 3f81b0... |                               | NULL | 1.0     |
| 3   | josephspurrier | success.sql | 2019-01-12 16:04:16 | 3             | 57cc0b... |                               | NULL | 1.0     |

## Migration File Specifications

There are a few components of a changeset:

- Header: must be prefixed by "--changeset " and must follow this format: `author:id` (single line, required)
- Body: valid sql text (multi-line, required)
- Description: must be prefixed by "--description " (multi-line, optional)
- Rollback: must be prefixed "--rollback "  (multi-line, optional)
- Include: must be prefixed by "--include " and must follow this format: `relativefilename.sql` (single line, optional)
- Comments: any other line that starts with "--" (multi-line, optional)

Blank lines are ignored by Rove. The prefixes above are strict so you cannot change the case or add spacing. For instance, you cannot add a space after the dashes: `-- changeset`.

Example migration file:

```sql
--changeset josephspurrier:1
--description Create the user status table.
--description Set deleted_at as a timestamp.
CREATE TABLE user_status (
    id TINYINT(1) UNSIGNED NOT NULL AUTO_INCREMENT,
    
    status VARCHAR(25) NOT NULL,
    
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
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

The rollback should be SQL which reverts the changes made by the changeset.

### Include

The include allows you to reference other changeset files to load. The filename should be a relative path.

### Comments

Any comments at the beginning of the lines are ignored. They do not count towards the checksum.
