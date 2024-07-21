package gateway

import "github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
