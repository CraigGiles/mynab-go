package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type System struct {
	database *sql.DB
	router   *mux.Router
}

func setup_routes(router *mux.Router) {
}

func initialize_system() System {
	var result System = System{}

	result.router = mux.NewRouter()
	setup_routes(result.router)

	result.database = create_database_connection("postgres", "root", "mynab") // TODO(craig) secrets and config

	return result
}

func create_database_connection(username string, password string, database string) *sql.DB {
	connection_string := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, database)

	sql.Open("postgres", connection_string)

	var db, err = sql.Open("postgres", connection_string)

	if err != nil {
		log.Fatal("Error opening connection to database: %s", err)
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Error with database connection: %s", err)
		panic(err)
	}

	return db
}

// Main function
func main() {
	system := initialize_system()

	// Start server
	fmt.Println("Listening on port ", 8080)
	log.Fatal(http.ListenAndServe(":8080", system.router))
}
