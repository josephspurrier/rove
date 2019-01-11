package rove_test

import (
	"log"

	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
)

func Example() {
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
	db, err := mysql.New(&mysql.Connection{
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
}
