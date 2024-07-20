package gateway

import "github.com/uiratan/fullcycle-archdev-microservices/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
