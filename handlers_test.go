package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func setupTestDB() {
	var err error
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic("Failed to open test database")
	}

	_, err = db.Exec(`
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);
	`)
	if err != nil {
		panic("Failed to create table")
	}

	// Insert test users
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", "testuser", hashedPassword, "user")
	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", "admin", hashedPassword, "admin")
	if err != nil {
		panic("Failed to insert test users")
	}
}

func TestLoginPage_SuccessfulLogin(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "testuser")
	reqBody.Set("password", "password123")

	req := httptest.NewRequest("POST", "/login", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	LoginPage(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 200 or 303, got %d", w.Code)
	}
}

func TestLoginPage_InvalidCredentials(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "testuser")
	reqBody.Set("password", "wrongpassword")

	req := httptest.NewRequest("POST", "/login", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	LoginPage(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestRegisterPage_SuccessfulRegistration(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "newuser")
	reqBody.Set("password", "newpassword")

	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	RegisterPage(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 Redirect, got %d", w.Code)
	}
}

func TestRegisterPage_DuplicateUser(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "testuser") // Already exists
	reqBody.Set("password", "newpassword")

	req := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	RegisterPage(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409 Conflict, got %d", w.Code)
	}
}

func TestAdminLogin_Successful(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "admin")
	reqBody.Set("password", "password123")

	req := httptest.NewRequest("POST", "/admin", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AdminPage(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 Redirect, got %d", w.Code)
	}
}

func TestAdminLogin_InvalidCredentials(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("username", "admin")
	reqBody.Set("password", "wrongpassword")

	req := httptest.NewRequest("POST", "/admin", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AdminPage(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestUpdateUser_Successful(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("id", "1")
	reqBody.Set("username", "updateduser")

	req := httptest.NewRequest("POST", "/admin/update", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	UpdateUser(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 Redirect, got %d", w.Code)
	}
}

func TestDeleteUser_Successful(t *testing.T) {
	setupTestDB()
	reqBody := url.Values{}
	reqBody.Set("id", "1")

	req := httptest.NewRequest("POST", "/admin/delete", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	DeleteUser(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 Redirect, got %d", w.Code)
	}
}
