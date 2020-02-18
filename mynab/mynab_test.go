package main

import (
	"testing"
	"time"
)

func TestAccountCreation(t *testing.T) {
	var account = make_account("chase", "checking")

	if account.name != "chase" {
		t.Error("Expected 'chase' but got ", account.name)
	}

	if len(account.transactions) != 0 {
		t.Error("Expected 0 transactions but got ", len(account.transactions))
	}
}

const (
	DateLayoutISO = "2006-01-02"
	DateLayoutUS  = "January 2, 2006"
)

func TestTransactionCreation(t *testing.T) {
	date, _ := time.Parse(DateLayoutISO, "2019-01-15")

	var transaction = make_transaction(date, "Bob Smith", "VENMO", 42.0)

	if transaction.date.Year() != 2019 ||
		transaction.date.Month() != 1 ||
		transaction.date.Day() != 15 {
		t.Error("Transaction date is incorrect: Expected ", date, ", got ", transaction.date)
	}

	t.Logf("Transaction(%s, %v/%v/%v, %s, %s, %d)",
		transaction.id,
		transaction.date.Year(), transaction.date.Day(), transaction.date.Month(),
		transaction.payee,
		transaction.category,
		transaction.amount)
}
