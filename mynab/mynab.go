package main

import (
	"github.com/google/uuid"
	"time"
)

type AccountType string

const (
	AccountType_Checking AccountType = "checking"
	AccountType_Savings  AccountType = "savings"
)

type Account struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Account_type AccountType   `json:"type"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	Payee    string    `json:"payee"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
}

func account_type_from_string(value string) AccountType {
	var result AccountType
	if value == "checking" {
		result = AccountType_Checking
	} else {
		result = AccountType_Savings
	}

	return result
}

func make_account(name string, account_type string) Account {
	var result Account

	id, _ := uuid.NewUUID()
	result.Id = id.String()

	result.Name = name
	result.Transactions = []Transaction{}
	result.Account_type = account_type_from_string(account_type)

	return result
}

func make_transaction(date time.Time, payee string, category string, amount float64) Transaction {
	var result Transaction

	id, _ := uuid.NewUUID()
	result.Id = id.String()

	result.Date = date
	result.Payee = payee
	result.Category = category
	result.Amount = amount

	return result
}
