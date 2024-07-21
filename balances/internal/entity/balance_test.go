package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBalance(t *testing.T) {
	balance := NewBalance("1", 100.00)
	assert.NotNil(t, balance)
	assert.Equal(t, "1", balance.AccountID)
	assert.Equal(t, 100.00, balance.Balance)
}
