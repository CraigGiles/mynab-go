package main

type AccountType int32

const (
	AccountType_Checking AccountType = 0
	AccountType_Savings
)

type Account struct {
	id           string      `json:"id"`
	account_type AccountType `json:"account_type"`
}
