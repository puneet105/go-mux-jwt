package models

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

type AccountTokenClaims struct{
	jwt.StandardClaims
	AccountNumber	int64	`json:"account_number"`
}
type Update struct {
	AccountNumber 	int64	`json:"account_number"`
	FirstName		string	`json:"first_name"`
	LastName		string	`json:"last_name"`
	Balance			int32	`json:"balance"`
}

type Login struct{
	AccountNumber	int64	`json:"account_number"`
	Password		string	`json:"password"`
}

type CreateAccount struct {
	FirstName	string		`json:"first_name"`
	LastName	string		`json:"last_name"`
	Password	string		`json:"pasword"`
}

type TransferRequest struct{
	AccountNumber	int64	`json:"account_number"`
	Amount			int32	`json:"amount"`
}

type Account struct{
	ID 				int			`json:"id"`
	FirstName		string		`json:"first_name"`
	LastName		string		`json:"last_name"`
	AccountNumber	int64		`json:"account_number"`
	EncryPassword	string		`json:"-"`
	Balance			int32		`json:"balance"`
	CreatedAt		time.Time	`json:"created_at"`
}

func NewAccount(firstname, lastname, password string)(*Account, error){
	encPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return nil,err
	}
	return &Account{
		FirstName: firstname,
		LastName: lastname,
		AccountNumber: int64(rand.Intn(1000000)),
		EncryPassword: string(encPassword),
		Balance: int32(rand.Intn(10000)),
		CreatedAt: time.Now().UTC(),
	},nil
}


func WriteJSON(w http.ResponseWriter, status int, message interface{})error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(message)
}