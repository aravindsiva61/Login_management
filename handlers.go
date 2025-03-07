package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword, role string
		err := db.QueryRow("SELECT password, role FROM users WHERE username = ?", username).Scan(&storedPassword, &role)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		} else if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			w.Write([]byte("User Login Successful!"))
		}
		return
	}
	templates.ExecuteTemplate(w, "login.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	templates.ExecuteTemplate(w, "register.html", nil)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ? AND role = 'admin'", username).Scan(&storedPassword)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "admin_login.html", nil)
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Fetch all non-admin users
	rows, _ := db.Query("SELECT id, username FROM users WHERE role != 'admin'")
	defer rows.Close()

	var users []struct {
		ID       int
		Username string
	}
	for rows.Next() {
		var u struct {
			ID       int
			Username string
		}
		rows.Scan(&u.ID, &u.Username)
		users = append(users, u)
	}

	templates.ExecuteTemplate(w, "admin.html", users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	newUsername := r.FormValue("username")

	if id == "" || newUsername == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE users SET username = ? WHERE id = ?", newUsername, id)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	_, _ = db.Exec("DELETE FROM users WHERE id=?", id)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
