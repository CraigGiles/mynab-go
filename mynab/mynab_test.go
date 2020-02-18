package main

import (
	"testing"
	"time"
)

func TestAccountCreation(t *testing.T) {
	var account = make_account("chase", "checking")

	if account.Name != "chase" {
		t.Error("Expected 'chase' but got ", account.Name)
	}

	if len(account.Transactions) != 0 {
		t.Error("Expected 0 transactions but got ", len(account.Transactions))
	}
}

const (
	DateLayoutISO = "2006-01-02"
	DateLayoutUS  = "January 2, 2006"
)

func TestTransactionCreation(t *testing.T) {
	date, _ := time.Parse(DateLayoutISO, "2019-01-15")

	var transaction = make_transaction(date, "Bob Smith", "VENMO", 42.0)

	if transaction.Date.Year() != 2019 ||
		transaction.Date.Month() != 1 ||
		transaction.Date.Day() != 15 {
		t.Error("Transaction date is incorrect: Expected ", date, ", got ", transaction.Date)
	}

	t.Logf("Transaction(%s, %v/%v/%v, %s, %s, %d)",
		transaction.Id,
		transaction.Date.Year(), transaction.Date.Day(), transaction.Date.Month(),
		transaction.Payee,
		transaction.Category,
		transaction.Amount)
}
