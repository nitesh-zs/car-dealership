package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

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
	id := vars["id"]

	// parse ID
	err := parseID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid ID"}}`)

		return
	}

	car, err := h.svc.GetByID(id)
	if err != nil {
		handleServerErr(err, id, w)
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

	// validate car
	err = validateCar(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":{"code":"invalid body","message":"%v"}}`, err)

		return
	}

	// validate params
	err = validateParams(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":{"code":"missing param(s)","requiredParams":"%v"}}`, err)

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

	// validate car
	err = validateCar(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":{"code":"invalid body","message":"%v"}}`, err)

		return
	}

	// validate params
	err = validateParams(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":{"code":"missing param(s)","requiredParams":"%v"}}`, err)

		return
	}

	id := mux.Vars(r)["id"]

	// parse ID
	err = parseID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid ID"}}`)

		return
	}

	car.ID = id

	updatedCar, err := h.svc.Update(&car)
	if err != nil {
		handleServerErr(err, id, w)
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
	id := mux.Vars(r)["id"]

	// parse ID
	err := parseID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":"invalid ID"}}`)

		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		handleServerErr(err, id, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return nil
}

func handleServerErr(err error, id string, w http.ResponseWriter) {
	if err == customErrors.CarNotExists() {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":"entity not found","id":"`+id+`"}}`)
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"code":"DB error"}}`)
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

func validateParams(car *model.Car) error {
	// validate necessary car parameters are passed
	missingParams := make([]string, 0)

	if car.Name == "" {
		missingParams = append(missingParams, model.ParamName)
	}

	if car.Brand == "" {
		missingParams = append(missingParams, model.ParamBrand)
	}

	if car.FuelType == "" {
		missingParams = append(missingParams, model.ParamFuelType)
	}

	if car.YearOfManufacture == 0 {
		missingParams = append(missingParams, model.ParamYearOfManufacture)
	}

	// validate engine params
	missingParams = append(missingParams, validateEngineParams(car)...)

	if len(missingParams) != 0 {
		return customErrors.MissingParams{
			RequiredParams: missingParams,
		}
	}

	return nil
}

func validateEngineParams(car *model.Car) []string {
	missingParams := make([]string, 0)

	switch car.FuelType {
	// for electric cars, range must be present
	case model.ValueElectric:
		if car.Engine.Range == 0 {
			missingParams = append(missingParams, model.ParamRange)
		}

	// for non-electric cars, displacement and noOfCylinders must be present
	case model.ValuePetrol, model.ValueDiesel:
		if car.Engine.Displacement == 0 {
			missingParams = append(missingParams, model.ParamDisplacement)
		}

		if car.Engine.NoOfCylinders == 0 {
			missingParams = append(missingParams, model.ParamNoOfCylinders)
		}
	}

	return missingParams
}

func validateCar(car *model.Car) error {
	// validate year of manufacture
	err := validateYearOfManufacture(car.YearOfManufacture)
	if err != nil {
		return err
	}

	// validate brand
	err = validateBrand(car.Brand)
	if err != nil {
		return err
	}

	// validate fuel type
	err = validateFuelType(car.FuelType)
	if err != nil {
		return err
	}

	return nil
}

func validateYearOfManufacture(year int) error {
	if (year < model.MinYear && year != 0) || year > time.Now().Year() {
		return customErrors.InvalidYOM()
	}

	return nil
}

func validateBrand(brand string) error {
	if brand != "" && brand != model.ValueTesla && brand != model.ValueFerrari && brand != model.ValuePorsche && brand != model.ValueBMW {
		return customErrors.InvalidBrand()
	}

	return nil
}

func validateFuelType(fuelType string) error {
	if fuelType != "" && fuelType != model.ValueElectric && fuelType != model.ValueDiesel && fuelType != model.ValuePetrol {
		return customErrors.InvalidFuelType()
	}

	return nil
}
