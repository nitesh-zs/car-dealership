package customerrors

import (
	"carAPI/model"
	"fmt"
)

type EntityNotExists string

type MissingParams struct {
	RequiredParams []string
}

type InvalidValue string

func (e EntityNotExists) Error() string {
	return fmt.Sprintf("%v not exists", string(e))
}

func (m MissingParams) Error() string {
	return fmt.Sprint(m.RequiredParams)
}

func (i InvalidValue) Error() string {
	return fmt.Sprintf("Invalid Value of %v", string(i))
}

func CarNotExists() EntityNotExists {
	var e EntityNotExists = "Car"
	return e
}

func EngineNotExists() EntityNotExists {
	var e EntityNotExists = "Engine"
	return e
}

func InvalidFuelType() InvalidValue {
	var e InvalidValue = model.ParamFuelType
	return e
}

func InvalidYOM() InvalidValue {
	var e InvalidValue = model.ParamYearOfManufacture
	return e
}

func InvalidBrand() InvalidValue {
	var e InvalidValue = model.ParamBrand
	return e
}
