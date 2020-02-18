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
	id           string        `json:"id"`
	name         string        `json:"name"`
	account_type AccountType   `json:"type"`
	transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	id       string    `json:"id"`
	date     time.Time `json:"date"`
	payee    string    `json:"payee"`
	category string    `json:"category"`
	amount   int64     `json:"amount"`
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
	result.id = id.String()

	result.name = name
	result.transactions = []Transaction{}
	result.account_type = account_type_from_string(account_type)

	return result
}

func make_transaction(date time.Time, payee string, category string, amount int64) Transaction {
	var result Transaction

	id, _ := uuid.NewUUID()
	result.id = id.String()

	result.date = date
	result.payee = payee
	result.category = category
	result.amount = amount

	return result
}
