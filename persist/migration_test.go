package persist_test

import (
	"database/sql"
	"testing"

	"github.com/intothevoid/kramerbot/persist"
	"go.uber.org/zap"
)

// Test migrateUserStoreFromJsonToDatabase() function
func TestMigrateUserStoreFromJsonToDatabase(t *testing.T) {
	// Create logger
	var logger = *zap.NewExample()

	err := persist.MigrateUserStoreFromJsonToDatabase(&logger)
	if err != nil {
		t.Errorf("Error migrating user store from json to database: %s", err)
	}

	// Open database file
	dbName := "user_test.db"
	defer DeleteDBFile(dbName)
	udb := persist.CreateDatabaseConnection(dbName, &logger)

	// Get count of users in database
	count, err := udb.DB.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		t.Errorf("Error getting count of users in database: %s", err)
	}

	actualCount := checkCount(count)

	if actualCount < 1 {
		t.Errorf("Expected at least 1 user in database, got %d", actualCount)
	}
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	return count
}
