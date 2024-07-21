package gateway

import "github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/entity"

type BalanceGateway interface {
	Save(balance *entity.Balance) error
	FindByID(id string) (*entity.Balance, error)
	FindByAccountID(accountID string) (*entity.Balance, error)
	Update(balance *entity.Balance) error
}
