package db

import "testing"

func TestDatabase(t *testing.T) {
	db := Database()

	if db == nil {
		t.Errorf("Expected non-nil database connection, got nil")
	}

	err := db.Ping()
	if err != nil {
		t.Errorf("Error pinging database: %v", err)
	}

	rows, err := db.Query("SELECT 1")
	if err != nil {
		t.Errorf("Error querying database: %v", err)
	}
	defer rows.Close()

	var result int
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			t.Errorf("Error scanning row: %v", err)
		}
	}

	if result != 1 {
		t.Errorf("Expected result to be 1, got %d", result)
	}
}
