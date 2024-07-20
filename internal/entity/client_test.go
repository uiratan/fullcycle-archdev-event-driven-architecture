package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewCliente(t *testing.T) {
	client, err := NewClient("Uiratan", "u@u.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "Uiratan", client.Name)
	assert.Equal(t, "u@u.com", client.Email)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("Uiratan", "u@u.com")
	err := client.Update("Liana", "l@l.com")
	assert.Nil(t, err)
	assert.Equal(t, "Liana", client.Name)
	assert.Equal(t, "l@l.com", client.Email)
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client, _ := NewClient("Uiratan", "u@u.com")
	err := client.Update("", "l@l.com")
	assert.Error(t, err, "name is required")
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("Uiratan", "u@u.com")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, len(client.Accounts), 1)
}
