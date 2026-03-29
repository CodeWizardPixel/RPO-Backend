package main

import (
	"fmt"
	"net/http"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Meow! Starting server...")

	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	fmt.Println("Database connection established!")

	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)
	mux.HandleFunc("/info", info)

	fmt.Println("Server on :8080")

	http.ListenAndServe(":8080", mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
}

func info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page\n")
}
