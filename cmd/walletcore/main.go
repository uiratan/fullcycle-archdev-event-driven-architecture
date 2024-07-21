package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/database"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/event"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_account"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_client"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/usecase/create_transaction"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/web"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/web/webserver"
	"github.com/uiratan/fullcycle-archdev-microservices/pkg/events"
)

func main() {
	// db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/wallet")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}

// func main() {
// 	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/wallet")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	fmt.Println("Connected to the database successfully!")
// }
