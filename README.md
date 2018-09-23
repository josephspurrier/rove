# rove

[![Go Report Card](https://goreportcard.com/badge/github.com/josephspurrier/rove)](https://goreportcard.com/report/github.com/josephspurrier/rove)
[![Build Status](https://travis-ci.org/josephspurrier/rove.svg)](https://travis-ci.org/josephspurrier/rove)
[![Coverage Status](https://coveralls.io/repos/github/josephspurrier/rove/badge.svg?branch=master&timestamp=20180923-01)](https://coveralls.io/github/josephspurrier/rove?branch=master)

## MySQL Database Migration Tool Based on Liquibase

This project is based off Liquibase, the database migration tool. It uses a slimmed down version of the DATABASECHANGELOG database table to store the applied changesets. It only supports SQL (no XML) migrations currently. For the most part, you can take your existing SQL migration files and use them with this tool. You can also import this package into your application and apply changesets without having to store the migrations on physical files. This allows you manage DB migrations inside of your application so you can distribute single binary applications.

**Note:** Do not run this tool on a database that already has Liquibase migrations applied - they are not compatiable because the way the checksums are calculated is different. The DATABASECHANGELOG database table which is used for storing the changesets is also different.

## Quick Start with Docker Compose

You can build a docker image from this repository and set it up along with a MySQL container using docker compose.

```bash
# Create a docker image.
docker build -t rove:latest .

# Launch MySQL and the rove with docker compose.
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