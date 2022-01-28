package main

import (
	"carAPI/handler"
	"carAPI/service"
	"carAPI/store"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// connecting to db
	db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println(err)
	}

	// initialize dependencies
	carStore := store.New(db)
	engineStore := store.NewEngineStore(db)
	svc := service.New(carStore, engineStore)
	h := handler.New(svc)

	// register handlers
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.HandleFunc("/car", h.HandleGetAll).Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", h.HandleGetByID).Methods(http.MethodGet)
	r.HandleFunc("/car", h.HandleCreate).Methods(http.MethodPost)
	r.HandleFunc("/car/{id}", h.HandleUpdate).Methods(http.MethodPut)
	r.HandleFunc("/car/{id}", h.HandleDelete).Methods(http.MethodDelete)

	// start server
	log.Println(http.ListenAndServe(":4000", r))
}
