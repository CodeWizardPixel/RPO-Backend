package main

import (
	"fmt"
	"go-back/handlers"
	"go-back/repository"
	"go-back/service"
	"net/http"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
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

	err = goose.SetDialect("sqlite3")
	if err != nil {
		fmt.Println("Error setting goose dialect:", err)
		return
	}

	err = goose.Up(db, "./data/migrations")
	if err != nil {
		fmt.Println("Error running migrations:", err)
		return
	}

	fmt.Println("Migrations completed successfully!")

	// terminalRepository := repository.NewTerminalRepository(db)
	// terminals, err := terminalRepository.GetAllTerminals()
	// if err != nil {
	// 	fmt.Println("Error fetching terminals:", err)
	// 	return
	// }
	// fmt.Println("Retrieved terminals:", len(terminals))

	UserRepository := repository.NewUserRepository(db)
	AuthService := service.NewAuthService(UserRepository, "your_secret_key")
	AuthHandler := handlers.NewAuthHandler(AuthService)

	
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/login", AuthHandler.GetToken)
	mux.HandleFunc("/api/v1/auth/validate", AuthHandler.ValidateToken)

	fmt.Println("Server on :8080")

	http.ListenAndServe(":8080", mux)
}