package main

import (
	"github.com/puneet105/go-mux-jwt/api"
	"github.com/puneet105/go-mux-jwt/database"
	"log"
)

var port = "9001"
func main(){
	conn, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewAPIServer(port, conn)
	err = server.RunServer()
	if err != nil{
		log.Println("Error starting server : ",err)
	}
}
