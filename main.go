package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"subscriptions/pkg/dbwork"
	"subscriptions/pkg/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/subscriptions", handlers.AddSubscriptions).Methods("POST")
	router.HandleFunc("/subscriptions/{id}", handlers.DeleteSubscriptions).Methods("DELETE")
	router.HandleFunc("/subscriptions/{id}", handlers.GetSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", handlers.UpdateSubscriptions).Methods("PUT")
	router.HandleFunc("/subscriptions", handlers.GetAllSubscriptions).Methods("GET")
	router.HandleFunc("/sum", handlers.CalculatePrice).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func init() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	configDB := dbwork.PostgresDBParams{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     port,
		SSLMode:  os.Getenv("SSLMODE"),
		DBName:   os.Getenv("DBNAME"),
	}

	fmt.Println(configDB)

	err = dbwork.InitializationPostgresDB(configDB)
	if err != nil {
		log.Fatal(err)
	}
}
