package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/puneet105/go-mux-jwt/models"
	"log"
)


type Storage interface {
	CreateAccount(*models.Account) error
	GetAccounts()([]*models.Account,error)
	DeleteAccount(int)error
	UpdateAccount(*models.Account, *models.Update)error
	GetAccountByID(int)(*models.Account,error)
	TransferToAccount(*models.Account, int32)error
	GetAccountByAccountNumber(int64)(*models.Account,error)
}
type PostgresStore struct{
	db *sql.DB
}

func NewPostgresConnection()(*PostgresStore, error){
	connStr := "user=postgres dbname=postgres password=gobank  sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return &PostgresStore{
		db: db,
	},nil
}

func (s *PostgresStore)Init()error{
	return s.CreateTable()
}

func (s *PostgresStore)CreateTable()error{
	query := `create table if not exists account(
		id serial primary key,
		first_name	varchar(50),
		last_name varchar(50),
		account_number serial,
    	encrypted_password varchar(100),
		balance serial,
    	created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore)CreateAccount(account *models.Account)error{
	query := `insert into account
	(first_name, last_name, account_number, encrypted_password, balance, created_at)
	values ($1,$2,$3,$4,$5,$6)`
	result, err := s.db.Query(query, account.FirstName, account.LastName, account.AccountNumber, account.EncryPassword, account.Balance, account.CreatedAt)
	if err != nil{
		return err
	}
	log.Printf("Result from the database is: %+v", result)
	return nil
}

func (s *PostgresStore)GetAccounts()(account []*models.Account,err error){
	rows, err := s.db.Query(`select * from account`)
	if err != nil{
		return nil, err
	}
	accounts := []*models.Account{}
	for rows.Next(){
		acc := new(models.Account)
		err := rows.Scan(&acc.ID,&acc.FirstName,&acc.LastName,&acc.AccountNumber,&acc.EncryPassword,&acc.Balance,&acc.CreatedAt)
		if err != nil{
			return  nil, err
		}
		accounts = append(accounts,acc)
	}
	return accounts,nil
}

func (s *PostgresStore)GetAccountByID(id int)(*models.Account, error){
	rows := s.db.QueryRow("select * from account where id =$1",id)
	account := new(models.Account)
		err := rows.Scan(&account.ID,&account.FirstName,&account.LastName,&account.AccountNumber,&account.EncryPassword,&account.Balance,&account.CreatedAt)
		if err != nil{
			return nil, err
		}
	return account,nil
}

func(s *PostgresStore)GetAccountByAccountNumber(accNumber int64)(*models.Account,error){
	rows := s.db.QueryRow("select * from account where account_number =$1",accNumber)
	account := new(models.Account)
	err := rows.Scan(&account.ID,&account.FirstName,&account.LastName,&account.AccountNumber,&account.EncryPassword,&account.Balance,&account.CreatedAt)
	if err != nil{
		return nil, err
	}
	return account,nil
}

func (s *PostgresStore)DeleteAccount(id int)error{
	_, err := s.db.Query("delete from account where id = $1",id)
	return err
}

func (s *PostgresStore)UpdateAccount(account *models.Account, update *models.Update)error{
	query := `update account set first_name = $1 , last_name = $2 , balance = $3 where account_number = $4 `
	result, err := s.db.Query(query, update.FirstName, update.LastName, update.Balance, account.AccountNumber)
	if err != nil{
		return err
	}
	log.Printf("Result from the database is : %+v", result)
	return nil
}

func(s *PostgresStore)TransferToAccount(account *models.Account, amount int32)error{
	query := `update account set balance = $1 where account_number = $2`
	result, err := s.db.Query(query, amount, account.AccountNumber)
	if err != nil{
		return err
	}
	log.Printf("Result from the database is : %+v", result)
	return nil
}
