package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	InitDB() // Initialize database

	r := mux.NewRouter()
	r.HandleFunc("/", LoginPage).Methods("GET", "POST")
	r.HandleFunc("/register", RegisterPage).Methods("GET", "POST")
	r.HandleFunc("/admin", AdminPage).Methods("GET", "POST")
	r.HandleFunc("/admin/dashboard", AdminDashboard).Methods("GET")
	r.HandleFunc("/admin/update", UpdateUser).Methods("POST")
	r.HandleFunc("/admin/delete", DeleteUser).Methods("POST")

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
