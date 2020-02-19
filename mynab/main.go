package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type System struct {
	database *sql.DB
	router   *mux.Router
}

func initialize_system() System {
	var result System = System{}

	result.router = mux.NewRouter()
	result.database = create_database_connection("postgres", "root", "mynab") // TODO(craig) secrets and config

	return result
}

func create_database_connection(username string, password string, database string) *sql.DB {
	connection_string := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, database)

	sql.Open("postgres", connection_string)

	var db, err = sql.Open("postgres", connection_string)

	if err != nil {
		// log.Fatal("Error opening connection to database: %s", err)
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		// log.Fatal("Error with database connection: %s", err)
		panic(err)
	}

	return db
}

//
//     -- Get Account --
// -----------------------------------------------------------------
func get_accounts_handler(s *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := s.database.Query("SELECT * FROM accounts")

		if err != nil {
			// TODO
		}

		var result []Account
		fmt.Println("Getting accounts:")

		for rows.Next() {
			var id string
			var name string
			var account_type string
			err = rows.Scan(&id, &name, &account_type)
			if err != nil {
				// TODO
			}

			var account Account
			account.Id = id
			account.Name = name
			account.Account_type = account_type_from_string(account_type)

			fmt.Printf("%v | %v | %v\n", account.Id, account.Name, account.Account_type)

			result = append(result, account)
		}

		json.NewEncoder(w).Encode(result)
	}
}

//
//     -- Add Account --
// -----------------------------------------------------------------
type AddAccountContext struct {
	Name         string `json:"name"`
	Account_type string `json:"type"`
}

func persist_account(db *sql.DB, account Account) bool {
	var lastInsertId int

	err := db.QueryRow("INSERT INTO accounts(id, name, type) VALUES($1,$2,$3) returning id;",
		account.Id, account.Name, account.Account_type).Scan(&lastInsertId)

	if err != nil {
		// TODO
	}

	fmt.Println("last inserted id =", lastInsertId)

	return true
}

func add_account_handler(s *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var ctx AddAccountContext
		json_err := json.NewDecoder(r.Body).Decode(&ctx)

		if json_err != nil {
			http.Error(w, json_err.Error(), http.StatusBadRequest)
		} else {
			account := make_account(ctx.Name, ctx.Account_type)
			_ = persist_account(s.database, account)

			json.NewEncoder(w).Encode(account)
		}
	}
}

//
//     -- Add Transaction --
// -----------------------------------------------------------------
const (
	DateLayoutISO = "2006-01-02"
	DateLayoutUS  = "January 2, 2006"
)

type AddTransactionContext struct {
	Date     string  `json:"date"`
	Payee    string  `json:"payee"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

func persist_transaction(db *sql.DB, t Transaction) bool {
	var last_insert_id int

	err := db.QueryRow("INSERT INTO transactions(id, date, payee, category, amount) VALUES($1,$2,$3,$4,$5) returning id;",
		t.Id,
		t.Date,
		t.Payee,
		t.Category,
		t.Amount).Scan(&last_insert_id)

	if err != nil {
		// TODO
	}

	fmt.Println("last inserted id =", last_insert_id)

	return true
}

func add_transaction_handler(s *System) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var ctx AddTransactionContext
		json_err := json.NewDecoder(r.Body).Decode(&ctx)

		if json_err != nil {
			http.Error(w, json_err.Error(), http.StatusBadRequest)
		} else {
			date, date_err := time.Parse(DateLayoutISO, ctx.Date)

			if date_err != nil {
				http.Error(w, date_err.Error(), http.StatusBadRequest)
			} else {
				transaction := make_transaction(date, ctx.Payee, ctx.Category, ctx.Amount)
				_ = persist_transaction(s.database, transaction)

				json.NewEncoder(w).Encode(transaction)
			}
		}
	}
}

//
//     -- Main System --
// -----------------------------------------------------------------
func setup_routes(system *System) {
	system.router.HandleFunc("/accounts", get_accounts_handler(system)).Methods("GET")
	system.router.HandleFunc("/accounts", add_account_handler(system)).Methods("POST")

	system.router.HandleFunc("/transactions", add_transaction_handler(system)).Methods("POST")

	// TODO add account link to a transaction

	// TODO get current account balance
	// TODO remove transaction
	// TODO update account balance
}

// Main function
func main() {
	system := initialize_system()
	setup_routes(&system)

	// Start server
	fmt.Println("Listening on port ", 8080)
	log.Fatal(http.ListenAndServe(":8080", system.router))
}
