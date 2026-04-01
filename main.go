// @title Lab2 REST API
// @version 1.0

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	httpSwagger "github.com/swaggo/http-swagger"

	_ "go-back/docs"
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
	AuthService := service.NewAuthService(UserRepository, "wruff")
	AuthHandler := handlers.NewAuthHandler(AuthService)

	TerminalRepository := repository.NewTerminalRepository(db)
	TerminalService := service.NewTerminalService(TerminalRepository, AuthService)
	TerminalHandler := handlers.NewTerminalHandler(TerminalService)

	KeyRepository := repository.NewKeyRepository(db)
	KeyService := service.NewKeyService(KeyRepository, AuthService)
	KeyHandler := handlers.NewKeyHandler(KeyService)

	CardRepository := repository.NewCardRepository(db)
	CardService := service.NewCardService(CardRepository, AuthService)
	CardHandler := handlers.NewCardHandler(CardService)

	TransactionRepository := repository.NewTransactionRepository(db)
	TransactionService := service.NewTransactionService(TransactionRepository, AuthService)
	TransactionHandler := handlers.NewTransactionHandler(TransactionService)

	UserService := service.NewUserService(UserRepository, AuthService)
	UserHandler := handlers.NewUserHandler(UserService)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/api/v1/auth/login", AuthHandler.GetToken)
	mux.HandleFunc("/api/v1/auth/validate", AuthHandler.ValidateToken)

	mux.HandleFunc("/api/v1/terminals/all", TerminalHandler.GetAllTerminals)
	mux.HandleFunc("/api/v1/terminals/get", TerminalHandler.GetTerminalByID)
	mux.HandleFunc("/api/v1/terminals/create", TerminalHandler.CreateTerminal)
	mux.HandleFunc("/api/v1/terminals/update", TerminalHandler.UpdateTerminal)
	mux.HandleFunc("/api/v1/terminals/delete", TerminalHandler.DeleteTerminal)

	mux.HandleFunc("/api/v1/users/all", UserHandler.GetAllUsers)
	mux.HandleFunc("/api/v1/users/get", UserHandler.GetUserByID)
	mux.HandleFunc("/api/v1/users/create", UserHandler.CreateUser)
	mux.HandleFunc("/api/v1/users/update", UserHandler.UpdateUser)
	mux.HandleFunc("/api/v1/users/delete", UserHandler.DeleteUser)

	mux.HandleFunc("/api/v1/keys/all", KeyHandler.GetAllKeys)
	mux.HandleFunc("/api/v1/keys/get", KeyHandler.GetKeyByID)
	mux.HandleFunc("/api/v1/keys/create", KeyHandler.CreateKey)
	mux.HandleFunc("/api/v1/keys/update", KeyHandler.UpdateKey)
	mux.HandleFunc("/api/v1/keys/delete", KeyHandler.DeleteKey)

	mux.HandleFunc("/api/v1/cards/all", CardHandler.GetAllCards)
	mux.HandleFunc("/api/v1/cards/get", CardHandler.GetCardByID)
	mux.HandleFunc("/api/v1/cards/create", CardHandler.CreateCard)
	mux.HandleFunc("/api/v1/cards/update", CardHandler.UpdateCard)
	mux.HandleFunc("/api/v1/cards/balance", CardHandler.UpdateCardBalance)
	mux.HandleFunc("/api/v1/cards/delete", CardHandler.DeleteCard)

	mux.HandleFunc("/api/v1/transactions/all", TransactionHandler.GetAllTransactions)
	mux.HandleFunc("/api/v1/transactions/get", TransactionHandler.GetTransactionByID)
	mux.HandleFunc("/api/v1/transactions/create", TransactionHandler.CreateTransaction)
	mux.HandleFunc("/api/v1/transactions/delete", TransactionHandler.DeleteTransaction)

	fmt.Println("Server on :8080")

	http.ListenAndServe(":8080", mux)
}