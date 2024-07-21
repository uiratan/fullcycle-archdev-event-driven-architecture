package create_balance

import (
	"time"

	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/entity"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/gateway"
)

type CreateBalanceInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type CreateBalanceOutputDTO struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewCreateBalanceUseCase(balanceGateway gateway.BalanceGateway) *CreateBalanceUseCase {
	return &CreateBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *CreateBalanceUseCase) Execute(input CreateBalanceInputDTO) (*CreateBalanceOutputDTO, error) {
	balance := entity.NewBalance(input.AccountID, input.Balance)

	err := uc.BalanceGateway.Save(balance)
	if err != nil {
		return nil, err
	}

	return &CreateBalanceOutputDTO{
		ID:        balance.ID,
		AccountID: balance.AccountID,
		Balance:   balance.Balance,
		CreatedAt: balance.CreatedAt,
	}, nil
}
