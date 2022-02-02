package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	customErrors "carAPI/custom-errors"
	"carAPI/model"
	"carAPI/service"
)

type handler struct {
	svc service.CarService
}

//nolint:revive //handler should not be exported
func New(s service.CarService) handler {
	return handler{svc: s}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	brand := q.Get("brand")
	withEngine := q.Get("withEngine")

	if withEngine == "" {
		withEngine = "false"
	}

	we, err := strconv.ParseBool(withEngine)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid value of withEngine","message":"withEngine must be true or false"}}`)

		return
	}

	cars, err := h.svc.GetAll(brand, we)
	if err != nil {
		handleServerErr(err, "", w)
		return
	}

	resp, err := json.Marshal(cars)
	if err != nil {
		handleMarshalErr(err, w)
	}

	_, _ = w.Write(resp)
}

func (h handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]

	car, err := h.svc.GetByID(ID)
	if err != nil {
		handleServerErr(err, ID, w)
		return
	}

	resp, err := json.Marshal(car)
	if err != nil {
		handleMarshalErr(err, w)
		return
	}

	_, _ = w.Write(resp)
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleParseErr(err, w)
		return
	}

	var car model.Car

	err = json.Unmarshal(body, &car)
	if err != nil {
		handleParseErr(err, w)
		return
	}

	// validate params
	if !validateParams(&car, w) {
		return
	}

	// validate car
	if !validateCar(&car, w) {
		return
	}

	newCar, err := h.svc.Create(&car)
	if err != nil {
		handleServerErr(err, "", w)
		return
	}

	resp, err := json.Marshal(newCar)
	if err != nil {
		handleMarshalErr(err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleParseErr(err, w)
		return
	}

	var car model.Car

	err = json.Unmarshal(body, &car)
	if err != nil {
		handleParseErr(err, w)
		return
	}

	// validate params
	if !validateParams(&car, w) {
		return
	}

	// validate car
	if !validateCar(&car, w) {
		return
	}

	id := mux.Vars(r)["id"]
	car.ID = id

	updatedCar, err := h.svc.Update(&car)
	if err != nil {
		handleServerErr(err, car.ID, w)
		return
	}

	resp, err := json.Marshal(updatedCar)
	if err != nil {
		handleMarshalErr(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]

	err := h.svc.Delete(ID)
	if err != nil {
		handleServerErr(err, ID, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleServerErr(err error, id string, w http.ResponseWriter) {
	switch err.(type) {
	default:
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"code":"DB error"}}`)

	case customErrors.EntityNotExists:
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":"entity not found","id":"`+id+`"}}`)
	}
}

func handleMarshalErr(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, `{"error":{"code":"Marshal error"}}`)
}

func handleParseErr(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, `{"error":{"code":"invalid body", "message":"cannot parse given body"}}`)
}

func validateParams(car *model.Car, w http.ResponseWriter) bool {
	// validate necessary car parameters are passed
	if car.Name == "" || car.Brand == "" || car.FuelType == "" || car.YearOfManufacture == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`)

		return false
	}

	// validate necessary engine parameters are passed
	if !validateEngine(car, w) {
		return false
	}

	return true
}

func validateCar(car *model.Car, w http.ResponseWriter) bool {
	// validate year of manufacture
	if !validateYearOfManufacture(car, w) {
		return false
	}

	// validate brand
	if !validateBrand(car, w) {
		return false
	}

	// validate fuel type
	if !validateFuelType(car, w) {
		return false
	}

	return true
}

func validateEngine(car *model.Car, w http.ResponseWriter) bool {
	switch car.FuelType {
	case "Electric":
		if !validateElectricEngine(car, w) {
			return false
		}

	default:
		if !validateNonElectricEngine(car, w) {
			return false
		}
	}

	return true
}

func validateYearOfManufacture(car *model.Car, w http.ResponseWriter) bool {
	if (car.YearOfManufacture < 1866 && car.YearOfManufacture != 0) || car.YearOfManufacture > time.Now().Year() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid body", "message":"invalid year of manufacture"}}`)

		return false
	}

	return true
}

func validateBrand(car *model.Car, w http.ResponseWriter) bool {
	if car.Brand != "" && car.Brand != "Tesla" && car.Brand != "Ferrari" && car.Brand != "Porsche" && car.Brand != "BMW" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid body", "message":"supported brands are Tesla, Ferrari, Porsche and BMW"}}`)

		return false
	}

	return true
}

func validateFuelType(car *model.Car, w http.ResponseWriter) bool {
	if car.FuelType != "" && car.FuelType != "Electric" && car.FuelType != "Diesel" && car.FuelType != "Petrol" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid body", "message":"fuelType must be Electric, Petrol or Diesel"}}`)

		return false
	}

	return true
}

func validateElectricEngine(car *model.Car, w http.ResponseWriter) bool {
	if car.Engine.Range == 0 || car.Engine.Displacement != 0 || car.Engine.NoOfCylinders != 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`)

		return false
	}

	return true
}

func validateNonElectricEngine(car *model.Car, w http.ResponseWriter) bool {
	if car.Engine.Range != 0 || car.Engine.Displacement == 0 || car.Engine.NoOfCylinders == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"missing param(s)", "requiredParams":["name", "yearOfManufacture","brand",
							"fuelType", "engine"],"engineParams":"either range or displacement and noOfCylinders must be passed"}}`)

		return false
	}

	return true
}
