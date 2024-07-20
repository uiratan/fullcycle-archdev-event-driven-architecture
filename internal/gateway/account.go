package gateway

import "github.com/uiratan/fullcycle-archdev-microservices/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindByID(id string) (*entity.Account, error)
}
