package service

import (
	"carAPI/model"
	"carAPI/store"
)

type service struct {
	carStore    store.CarStore
	engineStore store.EngineStore
}

//nolint:revive //service should not be exported
func New(c store.CarStore, e store.EngineStore) service {
	return service{
		carStore:    c,
		engineStore: e,
	}
}

func (s service) GetAll(brand string, withEngine bool) ([]model.Car, error) {
	cars, err := s.carStore.GetByBrand(brand)
	if err != nil {
		return nil, err
	}

	// if withEngine is false, then just return the cars
	if !withEngine {
		return cars, nil
	}

	// fetch all engines mapped to their IDs
	engines, err := s.engineStore.GetAll()
	if err != nil {
		return nil, err
	}

	// map engines to their corresponding cars
	for i := range cars {
		car := cars[i]
		eID := car.Engine.ID
		engine := engines[eID]

		cars[i].Engine = engine
	}

	return cars, nil
}

func (s service) GetByID(id string) (*model.Car, error) {
	car, err := s.carStore.GetByID(id)
	if err != nil {
		return nil, err
	}

	engine, err := s.engineStore.GetByID(car.Engine.ID)
	if err != nil {
		return nil, err
	}

	car.Engine = *engine

	return car, nil
}

func (s service) Create(car *model.Car) (*model.Car, error) {
	engine, err := s.engineStore.Create(&car.Engine)
	if err != nil {
		return nil, err
	}

	car.Engine.ID = engine.ID

	newCar, err := s.carStore.Create(car)
	if err != nil {
		return nil, err
	}

	newCar.Engine = *engine

	return newCar, nil
}

func (s service) Update(car *model.Car) (*model.Car, error) {
	carFromDB, err := s.carStore.GetByID(car.ID)
	if err != nil {
		return nil, err
	}

	updatedCar, err := s.carStore.Update(car)
	if err != nil {
		return nil, err
	}

	car.Engine.ID = carFromDB.Engine.ID

	updatedEngine, err := s.engineStore.Update(&car.Engine)
	if err != nil {
		return nil, err
	}

	updatedCar.Engine = *updatedEngine

	return updatedCar, nil
}

func (s service) Delete(id string) error {
	car, err := s.carStore.GetByID(id)
	if err != nil {
		return err
	}

	err = s.carStore.Delete(id)
	if err != nil {
		return err
	}

	err = s.engineStore.Delete(car.Engine.ID)
	if err != nil {
		return err
	}

	return nil
}
