package store

import "carAPI/model"

type CarStore interface {
	// GetByBrand gives all the cars of a given brand,
	// if empty string is passed as brand, then cars of all brands should be fetched
	GetByBrand(brand string) ([]model.Car, error)

	// GetByID fetches a car with given ID from DB
	GetByID(id string) (*model.Car, error)

	// Create creates a new car in DB
	Create(car *model.Car) (*model.Car, error)

	// Update updates an existing car in DB
	Update(car *model.Car) (*model.Car, error)

	// Delete deletes a car with given ID from DB
	Delete(id string) error
}

type EngineStore interface {
	// GetAll returns a mapping of all engines IDs to corresponding engines
	GetAll() (map[string]model.Engine, error)

	// GetByID fetches an engine with given ID from DB
	GetByID(id string) (*model.Engine, error)

	// Create creates a new engine in DB
	Create(engine *model.Engine) (*model.Engine, error)

	// Update updates an existing engine in DB
	Update(engine *model.Engine) (*model.Engine, error)

	// Delete deletes the engine with given ID from DB
	Delete(id string) error
}
