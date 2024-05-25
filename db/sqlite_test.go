package db_test

import (
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
