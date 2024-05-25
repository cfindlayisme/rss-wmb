package db

import (
	"database/sql"
	"log"

	"github.com/cfindlayisme/rss-wmb/env"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", env.GetStateFilePath())
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetIfLinkPrintedInDB(link string) bool {
	rows, err := db.Query("SELECT Link FROM FeedItems WHERE Link = ? AND Printed = 1", link)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			return false
		}
		return true
	}

	return false
}

func (db *DB) CleanDB() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	stmt, err := tx.Prepare("DELETE FROM FeedItems WHERE Timestamp < datetime('now', '-7 days')")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}

	res, err := stmt.Exec()
	if err != nil {
		log.Fatalf("Error executing statement: %v", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting rows affected: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Printf("Cleaned up %d items older than 7 days from the database", rowCount)
}

func (db *DB) WriteFeedItemsToDB(feedItemsNew map[string]bool) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS FeedItems (Link TEXT PRIMARY KEY, Printed BOOLEAN, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO FeedItems (Link, Printed) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	for link, printed := range feedItemsNew {
		_, err = stmt.Exec(link, printed)
		if err != nil {
			log.Fatalf("Error executing statement: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	for link := range feedItemsNew {
		log.Printf("Added %s to the database\n", link)
	}
}
