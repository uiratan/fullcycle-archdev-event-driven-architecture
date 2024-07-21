package main

import (
	"context"
	"database/sql"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/database"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/event"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/event/handler"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_account"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_client"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_transaction"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/web"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/web/webserver"
	"github.com/uiratan/fullcycle-archdev-microservices/pkg/events"
	"github.com/uiratan/fullcycle-archdev-microservices/pkg/kafka"
	"github.com/uiratan/fullcycle-archdev-microservices/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/wallet?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	// eventDispatcher := events.NewEventDispatcher()
	// eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	// eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	// transactionCreatedEvent := event.NewTransactionCreated()
	// balanceUpdatedEvent := event.NewBalanceUpdated()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))

	transactionCreatedEvent := event.NewTransactionCreated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}
