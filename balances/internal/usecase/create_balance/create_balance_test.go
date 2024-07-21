package create_balance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture-utils/balances/internal/usecase/mocks"
)

func TestCreateBalanceUseCase_Execute(t *testing.T) {
	balanceMock := &mocks.BalanceGatewayMock{}
	balanceMock.On("Save", mock.Anything).Return(nil)

	uc := NewCreateBalanceUseCase(balanceMock)
	inputDto := CreateBalanceInputDTO{
		AccountID: "c76a8e3b-21a7-439b-956f-cf37ee44d424",
		Balance:   100.00,
	}
	output, err := uc.Execute(inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	assert.Equal(t, output.AccountID, "c76a8e3b-21a7-439b-956f-cf37ee44d424")
	assert.Equal(t, output.Balance, 100.00)
	balanceMock.AssertExpectations(t)
	balanceMock.AssertNumberOfCalls(t, "Save", 1)
}
