package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/database"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/usecase/create_balance"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/usecase/find_account"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/web"
	"github.com/uiratan/fullcycle-archdev-event-driven-architecture/balances/internal/web/webserver"
)

type BalanceEvent struct {
	Name    string      `json:"Name"`
	Payload BalanceData `json:"Payload"`
}

type BalanceData struct {
	AccountIDFrom      string  `json:"account_id_from"`
	AccountIDTo        string  `json:"account_id_to"`
	BalanceAccountFrom float64 `json:"balance_account_id_from"`
	BalanceAccountTo   float64 `json:"balance_account_id_to"`
}

func main() {
	print("Starting server...")
	db, err := sql.Open("mysql", "root:root@tcp(mysql-balance:3306)/balances?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createTablesDb(db)
	populateDb(db)

	balanceDb := database.NewBalanceDB(db)

	createBalanceUseCase := create_balance.NewCreateBalanceUseCase(balanceDb)
	findAccountUseCase := find_account.NewFindAccountUseCase(balanceDb)

	go func() {
		webserver := webserver.NewWebServer(":3003")
		accountHandler := web.NewWebBalanceHandler(*findAccountUseCase)
		webserver.AddHandler("/balances/{account_id}", accountHandler.FindAccount)
		fmt.Println("Server running at port 3003")
		webserver.Start()
	}()

	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"client.id":         "balances",
		"group.id":          "balances",
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(configMap)
	if err != nil {
		fmt.Println("error consumer", err.Error())
	}

	topics := []string{"balances"}
	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value), msg.TopicPartition)

			var event BalanceEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Error JSON:", err)
				continue
			}

			input := create_balance.CreateBalanceInputDTO{
				AccountID: event.Payload.AccountIDFrom,
				Balance:   event.Payload.BalanceAccountFrom,
			}
			createBalanceUseCase.Execute(input)

			input = create_balance.CreateBalanceInputDTO{
				AccountID: event.Payload.AccountIDTo,
				Balance:   event.Payload.BalanceAccountTo,
			}
			createBalanceUseCase.Execute(input)
		}
		if msg != nil {
			c.CommitMessage(msg)
		}
	}

}

func createTablesDb(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS balances (id varchar(255), account_id varchar(255), balance integer, created_at timestamp)")
	if err != nil {
		panic(err)
	}
}

func populateDb(db *sql.DB) {
	_, err := db.Exec(`INSERT INTO balances (id, account_id, balance, created_at) VALUES 
										('3446e7a4-2a02-4c16-9cdf-62c0c31adfe4', '6b76d113-257a-4942-93fd-d97c42491dac', 1000, '2024-07-21'),
										('f3c897f4-2593-461e-a4a3-dd84b075dca0', 'cc457010-c5fe-4e94-9492-3e0aa0554625', 2000, '2024-07-21'),
										('f762d214-e880-4b73-bea5-317c2b454b4e', '2b42674d-10ce-49fb-a976-593d05198a42', 3000, '2024-07-21'),
										('61a2a5f4-e494-4f99-849b-5df122cd956d', 'b41337b8-13f5-413b-8112-3d221d6b15c1', 4000, '2024-07-21'),
										('860ce9fc-e243-4568-9ad7-03922fefc67d', '372b5874-252a-48ac-80d6-d9c4c2f10114', 5000, '2024-07-21')`)
	if err != nil {
		panic(err)
	}
}
