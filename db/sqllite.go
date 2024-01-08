package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetIfLinkPrintedInDB(link string) bool {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", os.Getenv("STATEFILE"))
	if err != nil {
		return false
	}
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT Link FROM FeedItems WHERE Link = ?", link)
	if err != nil {
		return false
	}
	defer rows.Close()

	// Iterate over the rows
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

func WriteFeedItemsToDB(feedItemsNew map[string]bool) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", os.Getenv("STATEFILE"))
	if err != nil {
		log.Fatalf("Error opening SQLite database: %v", err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS FeedItems (Link TEXT PRIMARY KEY, Printed BOOLEAN, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	// Prepare the statement
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO FeedItems (Link, Printed) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Insert the feed items into the table
	for link, printed := range feedItemsNew {
		_, err = stmt.Exec(link, printed)
		if err != nil {
			log.Fatalf("Error executing statement: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	// Print out the link of the feed items added
	for link := range feedItemsNew {
		log.Printf("Added %s to the database\n", link)
	}
}
