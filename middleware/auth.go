package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/puneet105/go-mux-jwt/database"
	"github.com/puneet105/go-mux-jwt/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

var secret = "puneetsharma105"

func validateToken(tokenString string)(bool, *models.AccountTokenClaims, error){
	token, err :=  jwt.ParseWithClaims(tokenString, &models.AccountTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil{
		return false, &models.AccountTokenClaims{}, errors.New("Error parsing token")
	}
	if claims, ok := token.Claims.(*models.AccountTokenClaims); ok && token.Valid{
		log.Printf("Claims are: %+v", claims)
		return true, claims, nil
	}
	return false, &models.AccountTokenClaims{}, errors.New("Error validating Token")
}

func CreateJwt(account *models.Account)(string, error){
	claims := models.AccountTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
		AccountNumber: account.AccountNumber,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func AuthMiddlewareHandler(handlerFunc http.HandlerFunc, s database.Storage)http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Inside Auth Middleare Handler Function...!!")
		tokenString := request.Header.Get("token")

		valid,tokenClaims, err := validateToken(tokenString)
		if err != nil{
			models.WriteJSON(writer, http.StatusForbidden, err)
			return
		}else if !valid{
			models.WriteJSON(writer, http.StatusForbidden, errors.New("Permission Denied...!!Invalid Token"))
			return
		}
		strId := mux.Vars(request)["id"]
		id, err := strconv.Atoi(strId)
		if err != nil{
			models.WriteJSON(writer, http.StatusBadRequest, errors.New("Invalid ID"))
			return
		}
		account, err := s.GetAccountByID(id)
		if err != nil{
			models.WriteJSON(writer, http.StatusForbidden, errors.New("Error fetching account by ID"))
			return
		}
		if account.AccountNumber != tokenClaims.AccountNumber{
			models.WriteJSON(writer, http.StatusForbidden, errors.New("Permission Denied...!!Account Not Found"))
			return
		}

		log.Println("Token is valid...!!")
		handlerFunc(writer,request)
	}

}
