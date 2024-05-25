package db_test

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/cfindlayisme/rss-wmb/db"
	"github.com/stretchr/testify/require"
)

func TestGetIfLinkPrintedInDB(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a query to select the link and return a row with the link
	link := "http://example.com"
	rows := sqlmock.NewRows([]string{"Link"}).AddRow(link)
	mock.ExpectQuery("SELECT Link FROM FeedItems WHERE Link = \\? AND Printed = 1").WithArgs(link).WillReturnRows(rows)

	// Call GetIfLinkPrintedInDB and check the result
	require.True(t, database.GetIfLinkPrintedInDB(link))

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanDB(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a transaction to be started
	mock.ExpectBegin()

	// Expect a query to delete the old feed items and return a result
	mock.ExpectPrepare("DELETE FROM FeedItems WHERE Timestamp < datetime\\('now', '-7 days'\\)")
	mock.ExpectExec("DELETE FROM FeedItems WHERE Timestamp < datetime\\('now', '-7 days'\\)").WillReturnResult(sqlmock.NewResult(0, 1))

	// Expect the transaction to be committed
	mock.ExpectCommit()

	// Call CleanDB
	database.CleanDB()

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWriteFeedItemsToDB(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a query to create the table
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS FeedItems \\(Link TEXT PRIMARY KEY, Printed BOOLEAN, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP\\)").WillReturnResult(sqlmock.NewResult(0, 0))

	// Expect a transaction to be started
	mock.ExpectBegin()

	// Expect a query to insert the feed items and return a result
	mock.ExpectPrepare("INSERT OR IGNORE INTO FeedItems \\(Link, Printed\\) VALUES \\(\\?, \\?\\)")
	link := "http://example.com"
	printed := true
	mock.ExpectExec("INSERT OR IGNORE INTO FeedItems \\(Link, Printed\\) VALUES \\(\\?, \\?\\)").WithArgs(link, printed).WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect the transaction to be committed
	mock.ExpectCommit()

	// Call WriteFeedItemsToDB with a map containing one feed item
	feedItems := map[string]bool{link: printed}
	database.WriteFeedItemsToDB(feedItems)

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNewDB(t *testing.T) {
	// Create a mock database
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Set the STATEFILE environment variable to a valid file path
	err = os.Setenv("STATEFILE", "/valid/path/to/db")
	require.NoError(t, err)

	// Call NewDB and check the result
	database, err := db.NewDB()
	require.NoError(t, err)
	require.NotNil(t, database)
}

func TestGetIfLinkPrintedInDB_Error(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a query to select the link and return an error
	link := "http://example.com"
	mock.ExpectQuery("SELECT Link FROM FeedItems WHERE Link = \\? AND Printed = 1").WithArgs(link).WillReturnError(errors.New("database error"))

	// Call GetIfLinkPrintedInDB and check the result
	require.False(t, database.GetIfLinkPrintedInDB(link))

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCleanDB_Error(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a transaction to begin and return an error
	mock.ExpectBegin().WillReturnError(errors.New("database error"))

	// Replace the standard logger's output with a custom writer
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Call CleanDB
	database.CleanDB()

	// Check the log output
	require.Contains(t, buf.String(), "Error beginning transaction: database error")

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWriteFeedItemsToDB_Error(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Create a DB instance with the mock database
	database := &db.DB{DB: mockDB}

	// Expect a query to create a table and return an error
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS FeedItems \\(Link TEXT PRIMARY KEY, Printed BOOLEAN, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP\\)").WillReturnError(errors.New("database error"))

	// Call WriteFeedItemsToDB and check the result
	feedItemsNew := map[string]bool{"http://example.com": true}
	err = database.WriteFeedItemsToDB(feedItemsNew)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database error")

	// Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
