package create_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/mocks"
)

func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &mocks.ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)

	uc := NewCreateClientUseCase(m)
	output, err := uc.Execute(CreateClientInputDTO{
		Name:  "Uiratan",
		Email: "u@u.com",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, output.Name, "Uiratan")
	assert.Equal(t, output.Email, "u@u.com")
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}
