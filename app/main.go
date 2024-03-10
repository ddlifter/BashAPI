package main

import (
	"log"
	"net/http"

	middleware "github.com/ddlifter/BashAPI/app/middleware"

	handlers "github.com/ddlifter/BashAPI/app/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/commands", handlers.GetCommands).Methods("GET")
	router.HandleFunc("/commands", handlers.AddCommand).Methods("POST")
	router.HandleFunc("/commands/{id}", handlers.GetCommand).Methods("GET")
	router.HandleFunc("/commands/{id}", handlers.DeleteCommand).Methods("DELETE")
	router.HandleFunc("/commands", handlers.DeleteCommands).Methods("DELETE")
	router.HandleFunc("/commands/run/{id}", handlers.RunCommand)
	router.HandleFunc("/commands/stop/{id}", handlers.StopCommand)
	router.Use(middleware.LoggingMiddleware)

	log.Fatal(http.ListenAndServe(":8000", router))
}
