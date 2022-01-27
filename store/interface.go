package store

import "carAPI/model"

type CarStore interface {
	// GetByBrand gives all the cars of a given brand,
	// if empty string is passed as brand, then cars of all brands should be fetched
	GetByBrand(brand string) ([]model.Car, error)

	GetByID(ID string) (model.Car, error)
	Create(car model.Car) (model.Car, error)
	Update(car model.Car) (model.Car, error)
	Delete(ID string) error
}

type EngineStore interface {
	GetByID(ID string) (model.Engine, error)
	Create(engine model.Engine) (model.Engine, error)
	Update(engine model.Engine) (model.Engine, error)
	Delete(ID string) error
}
