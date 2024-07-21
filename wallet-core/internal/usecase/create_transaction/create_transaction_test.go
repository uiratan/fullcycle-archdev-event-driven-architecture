package create_transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/entity"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/event"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/usecase/mocks"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/pkg/events"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("Uiratan", "u@u.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Liana", "l@l.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	dispatcher := events.NewEventDispatcher()
	eventTransactionCreated := event.NewTransactionCreated()
	eventBalanceUpdated := event.NewBalanceUpdated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransactionCreated, eventBalanceUpdated)
	output, err := uc.Execute(ctx, inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
