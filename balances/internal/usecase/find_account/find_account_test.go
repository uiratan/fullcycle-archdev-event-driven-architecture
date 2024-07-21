package find_account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/entity"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/usecase/mocks"
)

func TestFindAccountUseCase_Execute(t *testing.T) {
	balance := entity.NewBalance("c76a8e3b-21a7-439b-956f-cf37ee44d424", 1000.0)
	accountMock := &mocks.BalanceGatewayMock{}
	accountMock.On("FindByAccountID", balance.ID).Return(balance, nil)
	// balanceMock := &mocks.BalanceGatewayMock{}
	// balanceMock.On("Save", mock.Anything).Return(nil)

	uc := NewFindAccountUseCase(accountMock)
	inputDto := FindAccountInputDTO{
		AccountID: balance.ID,
	}

	output, err := uc.Execute(inputDto)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "FindByAccountID", 1)
	// balanceMock.AssertExpectations(t)
	// balanceMock.AssertNumberOfCalls(t, "Save", 1)
	assert.Nil(t, err)
	assert.Equal(t, balance.ID, output.ID)
	assert.Equal(t, balance.AccountID, output.AccountID)
	assert.Equal(t, balance.Balance, output.Balance)
	assert.Equal(t, balance.CreatedAt, output.CreatedAt)
}
