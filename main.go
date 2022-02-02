package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"carAPI/handler"
	"carAPI/middleware"
	"carAPI/service"
	"carAPI/store/car"
	"carAPI/store/engine"
)

func main() {
	// connecting to db
	db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	// initialize dependencies
	carStore := car.New(db)
	engineStore := engine.NewEngineStore(db)
	svc := service.New(carStore, engineStore)
	h := handler.New(svc)

	// register handlers
	r := mux.NewRouter()

	r.StrictSlash(true)

	r.HandleFunc("/car", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/car/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/car", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/car/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/car/{id}", h.Delete).Methods(http.MethodDelete)

	// set middlewares
	r.Use(middleware.AuthMiddleware)
	r.Use(middleware.RespHeaderMiddleware)

	// start server
	log.Println(http.ListenAndServe(":4000", r))
}
