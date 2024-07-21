package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/entity"
)

type BalanceGatewayMock struct {
	mock.Mock
}

func (m *BalanceGatewayMock) Save(balance *entity.Balance) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *BalanceGatewayMock) FindByAccountID(accountId string) (*entity.Balance, error) {
	args := m.Called(accountId)
	return args.Get(0).(*entity.Balance), args.Error(1)
}

// func (m *BalanceGatewayMock) FindByID(id string) (*entity.Balance, error) {
// 	args := m.Called(id)
// 	return args.Get(0).(*entity.Balance), args.Error(1)
// }

func (m *BalanceGatewayMock) Update(balance *entity.Balance) error {
	args := m.Called(balance)
	return args.Error(0)
}
