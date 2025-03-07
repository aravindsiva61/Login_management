package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestInitDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	InitDB()

	_, err = db.Query("SELECT id, username, role FROM users LIMIT 1")
	if err != nil {
		t.Fatalf("Table 'users' was not created properly: %v", err)
	}

	var username, role string
	err = db.QueryRow("SELECT username, role FROM users WHERE username = 'admin'").Scan(&username, &role)
	if err != nil {
		t.Fatalf("Admin user not found: %v", err)
	}
	if username != "admin" || role != "admin" {
		t.Errorf("Admin user data incorrect: got (%s, %s), expected (admin, admin)", username, role)
	}
}
