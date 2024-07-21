package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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
	fmt.Println("Creating tables...")
	fmt.Println("=============================")
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB: TABLE clients CREATED")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance integer, created_at date)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB: TABLE accounts CREATED")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount integer, created_at date)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB: TABLE transactions CREATED")
}

func populateDb(db *sql.DB) {
	fmt.Println("Populating DB...")
	var previousAccountID string
	for i := 1; i <= 5; i++ {
		fmt.Println("=============================")
		clientID := uuid.NewString()
		_, err := db.Exec(`INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)`,
			clientID, fmt.Sprintf("Client %d", i), fmt.Sprintf("client%d@example.com", i), time.Now())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Client inserted: ? ? ? ?", clientID, fmt.Sprintf("Client %d", i), fmt.Sprintf("client%d@example.com", i), time.Now())

		accountID := uuid.NewString()
		_, err = db.Exec(`INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)`,
			accountID, clientID, i*1000, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Account inserted: ? ? ? ?", accountID, clientID, i*1000, time.Now())

		if previousAccountID != "" {
			transactionID := uuid.NewString()
			_, err = db.Exec(`INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES (?, ?, ?, ?, ?)`,
				transactionID, previousAccountID, accountID, 100*i, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Transaction inserted: ? ? ? ?", transactionID, previousAccountID, accountID, 100*i, time.Now())
		}

		previousAccountID = accountID
	}
	fmt.Println("Records inserted successfully")
}
