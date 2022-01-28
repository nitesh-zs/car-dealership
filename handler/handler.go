package handler

import (
	customErrors "carAPI/custom-errors"
	"carAPI/model"
	"carAPI/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type handler struct {
	svc service.CarService
}

func New(s service.CarService) handler {
	return handler{svc: s}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check x-api-key in request header
		if r.Header.Get("x-api-key") != "nitesh-zs" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":{"code":"Authorization error","message":"A valid 'x-api-key' must be set in request headers"}}`))
			return
		}
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (h handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	brand := q.Get("brand")
	withEngine := q.Get("withEngine")

	var we bool
	if withEngine == "true" {
		we = true
	}

	cars, err := h.svc.GetAll(brand, we)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}

	resp, err := json.Marshal(cars)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}
	w.Write(resp)
}

func (h handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	ID := vars["id"]

	user, err := h.svc.GetByID(ID)
	if err != nil {
		switch err.(type) {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":{"code":"DB error"}}`))
			return

		case customErrors.CarNotExists:
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"entity not found","id":"` + ID + `"}}`))
			return
		}
	}

	resp, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}
	w.Write(resp)
}

func (h handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var car model.Car

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`))
		return
	}

	err = json.Unmarshal(body, &car)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`))
		return
	}

	// validate necessary car parameters are passed
	if car.Name == "" || car.Brand == "" || car.FuelType == "" || car.YearOfManufacture == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`))
		return
	}

	// validate necessary engine parameters are passed
	if !((car.Engine.Range != 0) || (car.Engine.Displacement != 0 && car.Engine.NoOfCylinders != 0)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`))
		return
	}

	// validate year of manufacture
	if car.YearOfManufacture < 1866 || car.YearOfManufacture > time.Now().Year() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"invalid year of manufacture"}}`))
		return
	}

	// validate brand
	if car.Brand != "Tesla" && car.Brand != "Ferrari" && car.Brand != "Porsche" && car.Brand != "BMW" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"supported brands are Tesla, Ferrari, Porsche and BMW"}}`))
		return
	}

	// validate fuel type
	if car.FuelType != "Electric" && car.FuelType != "Diesel" && car.FuelType != "Petrol" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"fuelType must be Electric, Petrol or Diesel"}}`))
		return
	}

	car, err = h.svc.Create(car)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}

	resp, err := json.Marshal(car)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var car model.Car

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`))
		return
	}

	err = json.Unmarshal(body, &car)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"cannot parse given body"}}`))
		return
	}

	// validate necessary car parameters are passed
	if car.Name == "" || car.Brand == "" || car.FuelType == "" || car.YearOfManufacture == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`))
		return
	}

	// validate necessary engine parameters are passed
	if !((car.Engine.Range != 0) || (car.Engine.Displacement != 0 && car.Engine.NoOfCylinders != 0)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`))
		return
	}

	// validate year of manufacture
	if car.YearOfManufacture < 1866 || car.YearOfManufacture > time.Now().Year() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"invalid year of manufacture"}}`))
		return
	}

	// validate brand
	if car.Brand != "Tesla" && car.Brand != "Ferrari" && car.Brand != "Porsche" && car.Brand != "BMW" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"supported brands are Tesla, Ferrari, Porsche and BMW"}}`))
		return
	}

	// validate fuel type
	if car.FuelType != "Electric" && car.FuelType != "Diesel" && car.FuelType != "Petrol" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":"invalid body", "message":"fuelType must be Electric, Petrol or Diesel"}}`))
		return
	}

	ID := mux.Vars(r)["id"]
	car.ID = ID

	car, err = h.svc.Update(car)
	if err != nil {
		switch err.(type) {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":{"code":"DB error"}}`))
			return

		case customErrors.CarNotExists:
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"entity not found","id":"` + ID + `"}}`))
			return
		}
	}

	resp, err := json.Marshal(car)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":{"code":"DB error"}}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}

func (h handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ID := mux.Vars(r)["id"]

	err := h.svc.Delete(ID)
	if err != nil {
		switch err.(type) {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":{"code":"DB error"}}`))
			return

		case customErrors.CarNotExists:
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":{"code":"entity not found","id":"` + ID + `"}}`))
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
