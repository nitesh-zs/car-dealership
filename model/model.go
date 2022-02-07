package model

type Engine struct {
	ID            string `json:"engineId"`
	Displacement  int    `json:"displacement"`
	NoOfCylinders int    `json:"noOfCylinders"`
	Range         int    `json:"range"`
}

type Car struct {
	ID                string `json:"carId"`
	Name              string `json:"name"`
	YearOfManufacture int    `json:"yearOfManufacture"`
	Brand             string `json:"brand"`
	FuelType          string `json:"fuelType"`
	Engine            Engine `json:"engine"`
}

const (
	ParamName              = "name"
	ParamYearOfManufacture = "yearOfManufacture"
	ParamBrand             = "brand"
	ParamFuelType          = "fuelType"
	ParamRange             = "range"
	ParamDisplacement      = "displacement"
	ParamNoOfCylinders     = "noOfCylinders"

	MinYear = 1866

	ValueElectric = "Electric"
	ValuePetrol   = "Petrol"
	ValueDiesel   = "Diesel"
	ValueTesla    = "Tesla"
	ValueFerrari  = "Ferrari"
	ValueBMW      = "BMW"
	ValuePorsche  = "Porsche"
)
