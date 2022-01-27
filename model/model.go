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
