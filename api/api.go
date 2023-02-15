package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/puneet105/go-mux-jwt/database"
	"github.com/puneet105/go-mux-jwt/middleware"
	"github.com/puneet105/go-mux-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct{
	port string
	store database.Storage
}
func NewAPIServer(port string, store database.Storage)*APIServer{
	return &APIServer{
		port: port,
		store: store,
	}
}
func (s *APIServer)RunServer()error{
	router := mux.NewRouter()
	router.HandleFunc("/account" , handleHTTPRoute(s.CreateAccountHandler))
	router.HandleFunc("/login", handleHTTPRoute(s.LoginHandler))
	router.HandleFunc("/account", middleware.AuthMiddlewareHandler(handleHTTPRoute(s.GetAccount), s.store))
	router.HandleFunc("/account/{id}", middleware.AuthMiddlewareHandler(handleHTTPRoute(s.AccountHandlerByID), s.store) )
	router.HandleFunc("/transfer", middleware.AuthMiddlewareHandler(handleHTTPRoute(s.TransferToAccount), s.store))
	router.HandleFunc("/update", middleware.AuthMiddlewareHandler(handleHTTPRoute(s.UpdateAccountHandler), s.store))
	log.Println("Server is listening on port : ",s.port)
	return http.ListenAndServe(":"+s.port, router)
}

func (s *APIServer)AccountHandlerByID(w http.ResponseWriter, r *http.Request)error{
	switch r.Method{
	case "GET":
		return s.GetAccountById(w,r)
	case "DELETE":
		return s.DeleteAccountById(w,r)
	default:
		return errors.New(fmt.Sprintf("method %s not allowed",r.Method))
	}
}

func(s *APIServer)LoginHandler(w http.ResponseWriter, r *http.Request)error{
	if r.Method != "POST"{
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
	var loginReq models.Login
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil{
		return err
	}
	defer r.Body.Close()
	account, err := s.store.GetAccountByAccountNumber(loginReq.AccountNumber)
	if err != nil{
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.EncryPassword), []byte(loginReq.Password)); err != nil{
		return fmt.Errorf("Authentication Failed!! Wrong Password")
	}
	token, err := middleware.CreateJwt(account)
	if err != nil{
		return err
	}

	log.Printf("Account Details are : %+v\n Token is : %s",account,token)

	return models.WriteJSON(w, http.StatusOK, fmt.Sprintf("Account with account number %d has successfully loggedIn\n" +
						"Token is : %s",account.AccountNumber,token))
}

func (s *APIServer)GetAccount(w http.ResponseWriter, r *http.Request) error{
	accounts, err := s.store.GetAccounts()
	if err != nil{
		return err
	}
	return models.WriteJSON(w,http.StatusOK,accounts)
}

func (s *APIServer)GetAccountById(w http.ResponseWriter, r *http.Request)error{
	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil{
		return fmt.Errorf("Invalid Id given %s", strId)
	}
	account, err := s.store.GetAccountByID(id)
	if err != nil{
		return err
	}
	return models.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer)CreateAccountHandler(w http.ResponseWriter, r *http.Request) error{
	var createAccReq models.CreateAccount
	if err := json.NewDecoder(r.Body).Decode(&createAccReq); err != nil{
		return err
	}
	defer r.Body.Close()
	account, err := models.NewAccount(createAccReq.FirstName, createAccReq.LastName, createAccReq.Password)
	if err != nil{
		return err
	}

	if err := s.store.CreateAccount(account); err != nil{
		return err
	}

	return models.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer)DeleteAccountById(w http.ResponseWriter, r *http.Request) error{
	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil{
		return fmt.Errorf("Invalid Id given %s", strId)
	}
	if err := s.store.DeleteAccount(id); err != nil{
		return err
	}
	return models.WriteJSON(w, http.StatusOK, fmt.Sprintf("Account with id %d has beeen deleted", id))
}

func (s *APIServer)TransferToAccount(w http.ResponseWriter, r *http.Request) error{
	transferReq := models.TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil  {
		return err
	}
	defer r.Body.Close()
	account, err := s.store.GetAccountByAccountNumber(transferReq.AccountNumber)
	if err != nil{
		return err
	}
	if err := s.store.TransferToAccount(account,transferReq.Amount); err != nil{
		return err
	}
	return models.WriteJSON(w, http.StatusOK, transferReq)
}

func(s *APIServer)UpdateAccountHandler(w http.ResponseWriter, r *http.Request)error{
	updateReq := models.Update{}
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil{
		return err
	}
	defer r.Body.Close()
	account, err := s.store.GetAccountByAccountNumber(updateReq.AccountNumber)
	if err != nil{
		return err
	}
	if err := s.store.UpdateAccount(account,&updateReq); err != nil{
		return err
	}
	return models.WriteJSON(w, http.StatusOK, updateReq)
}

type apiHandler func(http.ResponseWriter, *http.Request)error

func handleHTTPRoute(f apiHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err != nil{
			models.WriteJSON(w,http.StatusBadRequest, err.Error())
		}
	}
}