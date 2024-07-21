package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/uiratan/fullcycle-archdev-microservices/internal/entity"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
	clientFrom    *entity.Client
	clientTo      *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("Create table clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("Create table accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	db.Exec("Create table transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")
	s.transactionDB = NewTransactionDB(db)
	client1, err := entity.NewClient("Uiratan", "u@u.com")
	s.Nil(err)
	s.clientFrom = client1
	s.accountFrom = entity.NewAccount(s.clientFrom)
	s.accountFrom.Credit(1000)

	client2, err := entity.NewClient("Liana", "l@l.com")
	s.Nil(err)
	s.clientTo = client2
	s.accountTo = entity.NewAccount(s.clientTo)
	s.accountTo.Credit(1000)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 500)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)
}
