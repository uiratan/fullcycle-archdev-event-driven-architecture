package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("Uiratan", "u@u.com")
	account1 := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("Liana", "l@l.com")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, account1.Balance, 900.0)
	assert.Equal(t, account2.Balance, 1100.0)
}

func TestCreateTransactionWithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("Uiratan", "u@u.com")
	account1 := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("Liana", "l@l.com")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 2000)

	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, transaction)
	assert.Equal(t, account1.Balance, 1000.0)
	assert.Equal(t, account2.Balance, 1000.0)
}

func TestCreateTransactionWithZeroAmount(t *testing.T) {
	client1, _ := NewClient("Uiratan", "u@u.com")
	account1 := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("Liana", "l@l.com")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 0)

	assert.NotNil(t, err)
	assert.Error(t, err, "amount must be greater than zero")
	assert.Nil(t, transaction)
	assert.Equal(t, account1.Balance, 1000.0)
	assert.Equal(t, account2.Balance, 1000.0)
}

func TestCreateTransactionWithNegativeAmount(t *testing.T) {
	client1, _ := NewClient("Uiratan", "u@u.com")
	account1 := NewAccount(client1)
	account1.Credit(1000)
	client2, _ := NewClient("Liana", "l@l.com")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, -10)

	assert.NotNil(t, err)
	assert.Error(t, err, "amount must be greater than zero")
	assert.Nil(t, transaction)
	assert.Equal(t, account1.Balance, 1000.0)
	assert.Equal(t, account2.Balance, 1000.0)
}
