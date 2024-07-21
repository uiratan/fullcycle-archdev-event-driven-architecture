package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/database"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/event"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/event/handler"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/usecase/create_account"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/usecase/create_client"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/usecase/create_transaction"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/web"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/internal/web/webserver"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/pkg/events"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/pkg/kafka"
	"github.com/uiratan/fullcycle-archdev-microservices/wallet-core/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql-wallet:3306)/wallet?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createTablesDb(db)
	populateDb(db)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

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

	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}

func createTablesDb(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance integer, created_at date)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount integer, created_at date)`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB: TABLES CREATED")
}

func populateDb(db *sql.DB) {
	fmt.Println("Populating DB...")

	_, err := db.Exec(`INSERT INTO clients (id, name, email, created_at) VALUES 
											('15ca1636-1644-46a7-bd40-e75465c03d15', 'Client 1', 'client1@example.com', '2024-07-21'),
											('b788945b-e239-4b61-885c-fbac2674d9d8', 'Client 2', 'client2@example.com', '2024-07-21'),
											('4ae6d2c2-e6aa-4186-a216-34493946dc78', 'Client 3', 'client3@example.com', '2024-07-21'),
											('b39e08ca-f8d0-4ab4-aeaf-a77b22387b06', 'Client 4', 'client4@example.com', '2024-07-21'),
											('c365bcd3-b084-43d8-a028-805f8342b7d4', 'Client 5', 'client5@example.com', '2024-07-21')`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO accounts (id, client_id, balance, created_at) VALUES 
											('6b76d113-257a-4942-93fd-d97c42491dac', '15ca1636-1644-46a7-bd40-e75465c03d15', 1000, '2024-07-21'),
											('cc457010-c5fe-4e94-9492-3e0aa0554625', 'b788945b-e239-4b61-885c-fbac2674d9d8', 2000, '2024-07-21'),
											('2b42674d-10ce-49fb-a976-593d05198a42', '4ae6d2c2-e6aa-4186-a216-34493946dc78', 3000, '2024-07-21'),
											('b41337b8-13f5-413b-8112-3d221d6b15c1', 'b39e08ca-f8d0-4ab4-aeaf-a77b22387b06', 4000, '2024-07-21'),
											('372b5874-252a-48ac-80d6-d9c4c2f10114', 'c365bcd3-b084-43d8-a028-805f8342b7d4', 5000, '2024-07-21')`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES 
											('66cc50e7-ffb6-48e3-962d-af3740112a5c', '6b76d113-257a-4942-93fd-d97c42491dac', 'cc457010-c5fe-4e94-9492-3e0aa0554625', 200, '2024-07-21'),
											('9adea9fb-41ca-4fb9-8b45-838cb161bef7', 'cc457010-c5fe-4e94-9492-3e0aa0554625', '2b42674d-10ce-49fb-a976-593d05198a42', 300, '2024-07-21'),
											('b85a4d0a-c437-4d7d-82a7-5b241b2b7d52', '2b42674d-10ce-49fb-a976-593d05198a42', 'b41337b8-13f5-413b-8112-3d221d6b15c1', 400, '2024-07-21'),
											('83065b40-0e3f-4c24-84a3-e43d7d0b1e7e', 'b41337b8-13f5-413b-8112-3d221d6b15c1', '372b5874-252a-48ac-80d6-d9c4c2f10114', 500, '2024-07-21')`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB: POPULATED")
}
