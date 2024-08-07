package gateway

import "github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/entity"

type ClientGateway interface {
	Save(client *entity.Client) error
	Get(id string) (*entity.Client, error)
}
