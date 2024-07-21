package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBalance(t *testing.T) {
	balance := NewBalance("c76a8e3b-21a7-439b-956f-cf37ee44d424", 100.00)
	assert.NotNil(t, balance)
	assert.Equal(t, "c76a8e3b-21a7-439b-956f-cf37ee44d424", balance.AccountID)
	assert.Equal(t, 100.00, balance.Balance)
}
