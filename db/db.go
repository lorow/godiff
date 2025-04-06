package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite", "godiff.db")

	return err
}

func GetDBConnection() *sql.DB {
	return db
}

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS 'Project' (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'name' TEXT
	);`,
	`CREATE TABLE IF NOT EXISTS 'Request' (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'name' TEXT, 
		'url' TEXT, 
		'method' TEXT, 
		'headers' TEXT,
		'body' TEXT,  
		'response' INTEGER,
		'project_id' INTEGER,
		FOREIGN KEY(project_id) REFERENCES Project(id)
	);`,
	`CREATE TABLE IF NOT EXISTS 'Editor' (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'name' TEXT, 
		'content' TEXT, 
		'project_id' INTEGER,
		FOREIGN KEY(project_id) REFERENCES Project(id)
	);`,
}

type MigrationError struct {
	migration_id int
	err          error
}

func (e MigrationError) Error() string {
	return "Migration " + fmt.Sprint(e.migration_id) + " failed: " + e.err.Error()
}

func newMigrationError(id int, err error) MigrationError {
	return MigrationError{id, err}
}

func MigrateDatabase() error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for idx, migration := range migrations {
		_, err = tx.Exec(migration)

		if err != nil {
			tx.Rollback()
			return newMigrationError(idx, err)
		}
	}

	err = tx.Commit()
	return err
}
