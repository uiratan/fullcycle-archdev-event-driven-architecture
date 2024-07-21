package find_account

import (
	"time"

	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/gateway"
)

type FindAccountInputDTO struct {
	AccountID string `json:"account_id"`
}

type FindAccountOutputDTO struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type FindAccountUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewFindAccountUseCase(balanceGateway gateway.BalanceGateway) *FindAccountUseCase {
	return &FindAccountUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *FindAccountUseCase) Execute(input FindAccountInputDTO) (*FindAccountOutputDTO, error) {
	balance, err := uc.BalanceGateway.FindByAccountID(input.AccountID)
	if err != nil {
		return nil, err
	}

	return &FindAccountOutputDTO{
		ID:        balance.ID,
		AccountID: balance.AccountID,
		Balance:   balance.Balance,
		CreatedAt: balance.CreatedAt,
	}, nil
}
